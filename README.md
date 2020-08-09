# This pkg implements a non-thread safe ring buffer for buffering a byte stream

A ring buffer is useful in situations where you need a buffered reader. Once allocated a ring buffer does not need to allocate additional memory or move bigger junks of data around.

# Benchmarks
These are simple benchmarks which calculate an average throughput of data per second.

# BenchmarkChannelWithValueImpl  6184.24 MB/s
A ring buffer based on sending arrays of data. 

# BenchmarkChannelWithPtrImpl  249005.02 MB/s
A ring buffer where the written data is probably not copied because its throughput is incredibly high.

# BenchmarkSliceMovingImpl  11.85 MB/s
A ring buffer implementation based on moving data after reading data it.

# BenchmarkSliceWithAllocationImpl  1296.94 MB/s
A ring buffer based on reallocating memory.

# BenchmarkRingImpl  32184.91 MB/s
The ring buffer in this package.

# Conclusion 
If copying data is not needed when writing it, the channel implementation is unbeatable. If copying data on write is needed then this ring buffer implementations can probably be recommended.

# If you find obvious flaws here please let me know.
