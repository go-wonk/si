# Http Request

## Reusing http.Request with sync.Pool
I tried reusing http.Request with sync.Pool package, but it didn't make any big change. I brought some code from standard package to reset request body when getting from or putting back into the pool.

### Tests
- First test does not use any pool.
- Second and third tests reuses request body reader.
- Third test reuses response body reader
- Fourth test reuses request and response body reader.
```
BenchmarkHttpClient_DefaultPost-8   	                342	   4046098 ns/op	   33199 B/op	      88 allocs/op
BenchmarkHttpClient_DefaultPost_WithPool-8   	        370	   3142342 ns/op	   33095 B/op	      87 allocs/op
BenchmarkHttpClient_DefaultPost_WithPoolAndDoRead-8   	456	   2813064 ns/op	   20599 B/op	      81 allocs/op
BenchmarkHttpClient_RequestPost-8   	                400	   3210657 ns/op	   20221 B/op	      80 allocs/op

BenchmarkHttpClient_DefaultPost-8   	    			1555	    759026 ns/op	    5857 B/op	      74 allocs/op
BenchmarkHttpClient_DefaultPost_WithPool-8   	    	1348	    884274 ns/op	    5878 B/op	      74 allocs/op
BenchmarkHttpClient_DefaultPost_WithPoolAndDoRead-8   	1231	    981608 ns/op	    5478 B/op	      73 allocs/op
BenchmarkHttpClient_RequestPost-8   	    			1384	    733201 ns/op	    4858 B/op	      72 allocs/op
```
It is sufficient to reuse the request and response body. So, let's try not to tweak the standard package.

- The followings are test codes.
```go
const (
	testData = "a"
	// testDataRepeats = 1100
	testDataRepeats = 4096
	testUrl         = "https://127.0.0.1:8081/test/echo"
)

// BenchmarkHttpClient_DefaultPost-8   	     342	   4046098 ns/op	   33199 B/op	      88 allocs/op
func BenchmarkHttpClient_DefaultPost(b *testing.B) {
	if !onlinetest {
		b.Skip("skipping online tests")
	}
	siutils.AssertNotNilFailB(b, client)

	data := strings.Repeat(testData, testDataRepeats)
	url := testUrl

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {

		sendData := fmt.Sprintf("%s-%d", data, i)
		headerData := fmt.Sprintf("%d", i)

		request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(sendData)))
		siutils.AssertNilFailB(b, err)

		request.Header.Set("custom_header", headerData)

		resp, err := client.Do(request)
		siutils.AssertNilFailB(b, err)

		_, err = io.ReadAll(resp.Body)
		siutils.AssertNilFailB(b, err)

		resp.Body.Close()
	}
}

// BenchmarkHttpClient_DefaultPost_WithPool-8   	     370	   3142342 ns/op	   33095 B/op	      87 allocs/op
func BenchmarkHttpClient_DefaultPost_WithPool(b *testing.B) {
	if !onlinetest {
		b.Skip("skipping online tests")
	}
	siutils.AssertNotNilFailB(b, client)

	data := strings.Repeat(testData, testDataRepeats)
	url := testUrl

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {

		sendData := fmt.Sprintf("%s-%d", data, i)
		headerData := fmt.Sprintf("%d", i)

		buf := sicore.GetBytesReader([]byte(sendData))
		request, err := http.NewRequest(http.MethodPost, url, buf)
		siutils.AssertNilFailB(b, err)

		request.Header.Set("custom_header", headerData)

		resp, err := client.Do(request)
		siutils.AssertNilFailB(b, err)

		_, err = io.ReadAll(resp.Body)
		siutils.AssertNilFailB(b, err)

		resp.Body.Close()
		sicore.PutBytesReader(buf)
	}
}

// BenchmarkHttpClient_DefaultPost_WithPoolAndDoRead-8   	     456	   2813064 ns/op	   20599 B/op	      81 allocs/op
func BenchmarkHttpClient_DefaultPost_WithPoolAndDoRead(b *testing.B) {
	if !onlinetest {
		b.Skip("skipping online tests")
	}
	siutils.AssertNotNilFailB(b, client)

	data := strings.Repeat(testData, testDataRepeats)
	url := testUrl

	client := sihttp.NewHttpClient(client)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {

		sendData := fmt.Sprintf("%s-%d", data, i)
		headerData := fmt.Sprintf("%d", i)

		buf := sicore.GetBytesReader([]byte(sendData))
		request, err := http.NewRequest(http.MethodPost, url, buf)
		siutils.AssertNilFailB(b, err)

		request.Header.Set("custom_header", headerData)

		_, _, err = client.DoRead(request)
		siutils.AssertNilFailB(b, err)

		sicore.PutBytesReader(buf)
	}
}

// BenchmarkHttpClient_RequestPost-8   	     400	   3210657 ns/op	   20221 B/op	      80 allocs/op
func BenchmarkHttpClient_RequestPost(b *testing.B) {
	if !onlinetest {
		b.Skip("skipping online tests")
	}

	client := sihttp.NewHttpClient(client)

	data := strings.Repeat(testData, testDataRepeats)
	url := testUrl

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		sendData := fmt.Sprintf("%s-%d", data, i)
		headerData := fmt.Sprintf("%d", i)

		header := make(http.Header)
		header["custom_header"] = []string{headerData}
		client.RequestPost(url, header, []byte(sendData))
	}
}
```