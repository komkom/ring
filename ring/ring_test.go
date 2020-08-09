package ring

import (
	"bytes"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/bmizerany/assert"
	"github.com/stretchr/testify/require"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var letterRunes = []rune(letters)

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestRing_sequential(t *testing.T) {

	tests := []struct {
		ringSize  int
		writeSize int
		errWrite  error
		readSize  int
		loops     int
	}{
		{
			ringSize:  10,
			writeSize: 5,
			readSize:  5,
			loops:     100,
		},
		{
			ringSize:  11,
			writeSize: 5,
			readSize:  5,
			loops:     100,
		},
		{
			ringSize:  11,
			writeSize: 70,
			errWrite:  ErrOverflow,
			readSize:  5,
			loops:     100,
		},
	}

	for id, ts := range tests {

		t.Log(`idx`, id)

		r := New(ts.ringSize)

		cmpBuf := bytes.Buffer{}
		outBuf := bytes.Buffer{}

		readSlice := make([]byte, ts.readSize)

		for i := 0; i < ts.loops; i++ {

			seq := randSeq(ts.writeSize)
			length := r.Len()

			diff, err := r.Write([]byte(seq))
			if errors.Is(err, ErrOverflow) {

				for {
					n, err := r.Read(readSlice)
					require.NoError(t, err)
					diff -= n

					if n == 0 {
						break
					}

					outBuf.Write(readSlice[:n])
				}

				n, err := r.Write([]byte(seq))
				diff += n

				if ts.errWrite != nil {
					require.True(t, errors.Is(err, ts.errWrite))
					goto Next
				} else {
					require.NoError(t, err)
				}
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, length+diff, r.Len())
			cmpBuf.WriteString(seq)
		}

		for {
			length := r.Len()
			n, err := r.Read(readSlice)
			require.NoError(t, err)

			assert.Equal(t, length-n, r.Len())

			if n == 0 {
				break
			}

			outBuf.Write(readSlice[:n])
		}

		assert.Equal(t, cmpBuf.Bytes(), outBuf.Bytes())
	Next:
	}
}

const MB1 = 1000 * 1000

func BenchmarkChannelWithValueImpl(b *testing.B) {

	b.StopTimer()

	const size = 1000
	ch := make(chan [size]byte, MB1*2)

	tmp := []byte(randSeq(size))
	data := [size]byte{}
	copy(data[:], tmp)

	time.Sleep(time.Second)
	b.StartTimer()

	var counter int
	for i := 0; i < b.N; i++ {

		if counter < 1000 {
			for counter < MB1 {

				ch <- data
				counter += len(data)
			}
		}

		d := <-ch
		counter -= len(d)

		b.SetBytes(int64(len(d)))
	}

}

func BenchmarkChannelWithPtrImpl(b *testing.B) {

	b.StopTimer()

	size := 10000
	data := []byte(randSeq(size))
	ch := make(chan []byte, MB1*2)

	time.Sleep(time.Second)
	b.StartTimer()

	var counter int
	for i := 0; i < b.N; i++ {

		if counter < 1000 {
			for counter < MB1 {

				ch <- data
				counter += len(data)
			}
		}

		d := <-ch
		counter -= len(d)

		b.SetBytes(int64(len(d)))
	}
}

func BenchmarkSliceMovingImpl(b *testing.B) {

	b.StopTimer()

	buf := make([]byte, 0, MB1*5)
	data := []byte(letters)
	readBuf := make([]byte, 64)

	time.Sleep(time.Second)
	b.StartTimer()

	for i := 0; i < b.N; i++ {

		if len(buf) <= 1000 {
			for len(buf) < MB1 {
				buf = append(buf, data...)
			}
		}

		copy(readBuf, buf)
		buf = append(buf[:0], buf[64:]...)

		b.SetBytes(int64(len(readBuf)))
	}
}

func BenchmarkSliceWithAllocationImpl(b *testing.B) {

	b.StopTimer()

	buf := make([]byte, 0, MB1*5)
	data := []byte(letters)
	readBuf := make([]byte, 64)

	time.Sleep(time.Second)
	b.StartTimer()

	for i := 0; i < b.N; i++ {

		if len(buf) <= 1000 {
			for len(buf) < MB1 {
				buf = append(buf, data...)
			}
		}

		copy(readBuf, buf)
		buf = buf[64:]

		b.SetBytes(int64(len(readBuf)))
	}
}

func BenchmarkRingImpl(b *testing.B) {

	b.StopTimer()

	r := New(MB1 * 5)
	size := 10000
	data := []byte(randSeq(size))
	readS := make([]byte, size)

	b.StartTimer()

	for i := 0; i < b.N; i++ {

		if r.Len() <= 1000 {
			for r.Len() < MB1 {

				_, err := r.Write(data)
				if err != nil {
					panic(`write failed`)
				}
			}
		}

		n, err := r.Read(readS)
		if err != nil {
			panic(`read failed`)
		}

		b.SetBytes(int64(n))
	}
}
