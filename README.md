# This pkg implements a non-thread safe ring buffer for buffering a byte stream

A ring buffer can be useful in situations where you have a writer an a reader and you want to keep a bigger junk of memory (>1MB) in ram. Once allocated a ring buffer does not need to allocate additional memory or move bigger junks of data around.
As seen below this ring buffer has reasonable performance but a ring buffer implementation with a buffered chan in go is still faster.

# Benchmarks
These are simple benchmarks which calculate an average throughput of data per second.

## BenchmarkChannelImpl-12 	5834.70 MB/s
A go channel base ring buffer.

## BenchmarkSliceMovingImpl-12  11.89 MB/s
A ring buffer implementation based on moving data when reading data.

## BenchmarkSliceWithAllocationImpl-12  1267.92 MB/s
A ring buffer based on reallocating memory.

## BenchmarkRingImpl-12	 4193.45 MB/s
This ring buffer implementation.

## The winner go buffered channels
These benchmarks seem to indicate that a go buffered channel based ring buffer is the best performing.

# If you find obvious flaws here please let me know.
