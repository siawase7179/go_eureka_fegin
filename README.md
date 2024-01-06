# eureka-fegin

eureka와 instance 통신은 구현이 많이 되어 있지만, instance간 통신은 패키지가 따로 없어 구현해 보았다.

```go
module github.com/siawase7179/go_eureka_fegin

go 1.19

require (
	github.com/ArthurHlt/go-eureka-client v1.1.0
	github.com/sirupsen/logrus v1.9.3
)

require golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect

```

ArthurHit/go-ureka-client를 사용하였다.

```go
const (
	serviceName = "GO-SERVICE"
)

func init() {
	eureka.Init([]string{"http://localhost:8761/eureka"})

	eureka.NewInstance(serviceName, "localhost", 8082)
}
```
> [!note]
> 서비스 등록 / HearBeat /삭제 는 ArthurHit/go-ureka-client 에서 제공하는 방식과 동일하다.

feign.go
```go
app, err := eureka.GetApplication(serviceName)
	if err != nil {
		t.Fatal(err.Error())
	}
	feign.Append(*app)
```
eureka로부터 Application 정보를 가져와 feign에 Append 해준다.

```go
response, err := feign.Request(serviceName, feign.RequeustOption{
		Method: "GET",
	})
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(response.Response.Status)
```
이후 feign 패키지의 Request함수와 RequestOption 구조체로 라운드 로빈으로 Eureka Application으로 요청한다.

```go
type RequeustOption struct {
	Method string
	Path   string
	Body   string
	Header map[string]string
}

type FeignResponse struct {
	Response http.Response
	Body     []byte
}
```

> [!note]
> 아직 Instance의 연결 실패시 라운드로빈 타겟에서 제거가 되지 않는다.
>
> 추가 구현이 필요하다.
