# array-map

- concurrent
- append-only
- fixed-capacity

for real-time stream data processing

We use this pkg to process stock quotation ticks in financial scenes.

## Usage

```go
    // an array-map<string, int> with 100 keys capacity 
    //default key type is string
    arrmap.NewArrayMap[int](100)

    // an array-map<int32, string>
    // need a hash function for your custom key
    arrmap.NewArrayMapWithHasher[int32, string](100, func(k int32) uint32 {
        if k < 0 {
            return uint32(-k)
        }
        return uint32(k)
    })


```

## Benchmark

```shell
BenchmarkMultiGetSetDifferent
BenchmarkMultiGetSetDifferent-8          	 2394134	       493.1 ns/op
BenchmarkMultiGetSetDifferentSyncMap
BenchmarkMultiGetSetDifferentSyncMap-8   	  299953	      4722 ns/op
BenchmarkMultiGetSetDifferentCMap
BenchmarkMultiGetSetDifferentCMap-8      	 1000000	      2411 ns/op

BenchmarkMultiGetSetBlock
BenchmarkMultiGetSetBlock-8              	 2860682	       431.6 ns/op
BenchmarkMultiGetSetBlockSyncMap
BenchmarkMultiGetSetBlockSyncMap-8       	 1626556	       770.5 ns/op
BenchmarkMultiGetSetBlockCMap
BenchmarkMultiGetSetBlockCMap-8          	 1000000	      2312 ns/op
```
the `GetSetBlock` operation involves multiple read and write actions for the same set of keys, with the read and write frequencies being identical. 

This scenario closely resembles real-time stream processing data read and write operations.

