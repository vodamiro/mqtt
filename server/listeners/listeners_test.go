package listeners

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	l := New(nil)
	require.NotNil(t, l.internal)
}

func BenchmarkNewListeners(b *testing.B) {
	for n := 0; n < b.N; n++ {
		New(nil)
	}
}

func TestAddListener(t *testing.T) {
	l := New(nil)
	l.Add(NewMockListener("t1", ":1882"))
	require.Contains(t, l.internal, "t1")
}

func BenchmarkAddListener(b *testing.B) {
	l := New(nil)
	mocked := NewMockListener("t1", ":1882")
	for n := 0; n < b.N; n++ {
		l.Add(mocked)
	}
}

func TestGetListener(t *testing.T) {
	l := New(nil)
	l.Add(NewMockListener("t1", ":1882"))
	l.Add(NewMockListener("t2", ":1882"))
	require.Contains(t, l.internal, "t1")
	require.Contains(t, l.internal, "t2")

	g, ok := l.Get("t1")
	require.Equal(t, true, ok)
	require.Equal(t, g.ID(), "t1")
}

func BenchmarkGetListener(b *testing.B) {
	l := New(nil)
	l.Add(NewMockListener("t1", ":1882"))
	for n := 0; n < b.N; n++ {
		l.Get("t1")
	}
}

func TestLenListener(t *testing.T) {
	l := New(nil)
	l.Add(NewMockListener("t1", ":1882"))
	l.Add(NewMockListener("t2", ":1882"))
	require.Contains(t, l.internal, "t1")
	require.Contains(t, l.internal, "t2")
	require.Equal(t, 2, l.Len())
}

func BenchmarkLenListener(b *testing.B) {
	l := New(nil)
	l.Add(NewMockListener("t1", ":1882"))
	for n := 0; n < b.N; n++ {
		l.Len()
	}
}

func TestDeleteListener(t *testing.T) {
	l := New(nil)
	l.Add(NewMockListener("t1", ":1882"))
	require.Contains(t, l.internal, "t1")

	l.Delete("t1")
	_, ok := l.Get("t1")
	require.Equal(t, false, ok)
	require.Nil(t, l.internal["t1"])
}

func BenchmarkDeleteListener(b *testing.B) {
	l := New(nil)
	l.Add(NewMockListener("t1", ":1882"))
	for n := 0; n < b.N; n++ {
		l.Delete("t1")
	}
}

func TestServeListener(t *testing.T) {
	l := New(nil)
	l.Add(NewMockListener("t1", ":1882"))
	l.Serve("t1", MockEstablisher)
	time.Sleep(time.Millisecond)
	require.Equal(t, true, l.internal["t1"].(*MockListener).IsServing)

	l.Close("t1", MockCloser)
	require.Equal(t, false, l.internal["t1"].(*MockListener).IsServing)
}

func BenchmarkServeListener(b *testing.B) {
	l := New(nil)
	l.Add(NewMockListener("t1", ":1882"))
	for n := 0; n < b.N; n++ {
		l.Serve("t1", MockEstablisher)
	}
}

func TestServeAllListeners(t *testing.T) {
	l := New(nil)
	l.Add(NewMockListener("t1", ":1882"))
	l.Add(NewMockListener("t2", ":1882"))
	l.Add(NewMockListener("t3", ":1882"))
	l.ServeAll(MockEstablisher)
	time.Sleep(time.Millisecond)

	require.Equal(t, true, l.internal["t1"].(*MockListener).IsServing)
	require.Equal(t, true, l.internal["t2"].(*MockListener).IsServing)
	require.Equal(t, true, l.internal["t3"].(*MockListener).IsServing)

	l.Close("t1", MockCloser)
	l.Close("t2", MockCloser)
	l.Close("t3", MockCloser)

	require.Equal(t, false, l.internal["t1"].(*MockListener).IsServing)
	require.Equal(t, false, l.internal["t2"].(*MockListener).IsServing)
	require.Equal(t, false, l.internal["t3"].(*MockListener).IsServing)
}

func BenchmarkServeAllListeners(b *testing.B) {
	l := New(nil)
	l.Add(NewMockListener("t1", ":1882"))
	l.Add(NewMockListener("t2", ":1883"))
	l.Add(NewMockListener("t3", ":1884"))
	for n := 0; n < b.N; n++ {
		l.ServeAll(MockEstablisher)
	}
}

func TestCloseListener(t *testing.T) {
	l := New(nil)
	mocked := NewMockListener("t1", ":1882")
	l.Add(mocked)
	l.Serve("t1", MockEstablisher)
	time.Sleep(time.Millisecond)
	var closed bool
	l.Close("t1", func(id string) {
		closed = true
	})
	require.Equal(t, true, closed)
}

func BenchmarkCloseListener(b *testing.B) {
	l := New(nil)
	mocked := NewMockListener("t1", ":1882")
	l.Add(mocked)
	l.Serve("t1", MockEstablisher)
	for n := 0; n < b.N; n++ {
		l.internal["t1"].(*MockListener).done = make(chan bool)
		l.Close("t1", MockCloser)
	}
}

func TestCloseAllListeners(t *testing.T) {
	l := New(nil)
	l.Add(NewMockListener("t1", ":1882"))
	l.Add(NewMockListener("t2", ":1882"))
	l.Add(NewMockListener("t3", ":1882"))
	l.ServeAll(MockEstablisher)
	time.Sleep(time.Millisecond)
	require.Equal(t, true, l.internal["t1"].(*MockListener).IsServing)
	require.Equal(t, true, l.internal["t2"].(*MockListener).IsServing)
	require.Equal(t, true, l.internal["t3"].(*MockListener).IsServing)

	closed := make(map[string]bool)
	l.CloseAll(func(id string) {
		closed[id] = true
	})
	require.Contains(t, closed, "t1")
	require.Contains(t, closed, "t2")
	require.Contains(t, closed, "t3")
	require.Equal(t, true, closed["t1"])
	require.Equal(t, true, closed["t2"])
	require.Equal(t, true, closed["t3"])
}

func BenchmarkCloseAllListeners(b *testing.B) {
	l := New(nil)
	mocked := NewMockListener("t1", ":1882")
	l.Add(mocked)
	l.Serve("t1", MockEstablisher)
	for n := 0; n < b.N; n++ {
		l.internal["t1"].(*MockListener).done = make(chan bool)
		l.Close("t1", MockCloser)
	}
}
