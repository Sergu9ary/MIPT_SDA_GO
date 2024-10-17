//go:build !solution

//package jsonrpc
//
//import (
//	"bytes"
//	"context"
//	"encoding/json"
//	"errors"
//	"io"
//	"net/http"
//	"reflect"
//	"strings"
//)
//
//func MakeHandler(service interface{}) http.Handler {
//	serviceVal := reflect.ValueOf(service)
//	serviceType := serviceVal.Type()
//
//	methods := make(map[string]int)
//	for i := 0; i < serviceType.NumMethod(); i++ {
//		method := serviceType.Method(i)
//		methods[method.Name] = i
//	}
//
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		parts := strings.Split(r.URL.Path, "/")
//		methodName := parts[len(parts)-1]
//		methodIdx, ok := methods[methodName]
//		if !ok {
//			http.Error(w, "Method not found", http.StatusNotFound)
//			return
//		}
//
//		method := serviceType.Method(methodIdx)
//		reqType := method.Type.In(2).Elem()
//		reqVal := reflect.New(reqType)
//		if err := json.NewDecoder(r.Body).Decode(reqVal.Interface()); err != nil {
//			http.Error(w, "Invalid request format", http.StatusBadRequest)
//			return
//		}
//
//		results := method.Func.Call([]reflect.Value{
//			serviceVal,
//			reflect.ValueOf(context.Background()),
//			reqVal,
//		})
//
//		if errVal := results[1].Interface(); errVal != nil {
//			http.Error(w, errVal.(error).Error(), http.StatusInternalServerError)
//			return
//		}
//
//		rspVal := results[0].Interface()
//		w.Header().Set("Content-Type", "application/json")
//		if err := json.NewEncoder(w).Encode(rspVal); err != nil {
//			http.Error(w, "Response serialization failed", http.StatusInternalServerError)
//			return
//		}
//	})
//}
//
//func Call(ctx context.Context, endpoint string, method string, req, rsp interface{}) error {
//	reqData, err := json.Marshal(req)
//	if err != nil {
//		return err
//	}
//	url := endpoint + "/" + method
//	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqData))
//	if err != nil {
//		return err
//	}
//	httpReq.Header.Set("Content-Type", "application/json")
//
//	client := &http.Client{}
//	httpResp, err := client.Do(httpReq)
//	if err != nil {
//		return err
//	}
//	defer httpResp.Body.Close()
//	if httpResp.StatusCode != http.StatusOK {
//		body, _ := io.ReadAll(httpResp.Body)
//		return errors.New("RPC error: " + string(body))
//	}
//	return json.NewDecoder(httpResp.Body).Decode(rsp)
//}

package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
)

// Registry to map method names to their indices
func createMethodRegistry(serviceType reflect.Type) map[string]int {
	registry := make(map[string]int)
	for i := 0; i < serviceType.NumMethod(); i++ {
		method := serviceType.Method(i)
		registry[method.Name] = i
	}
	return registry
}

// Parse method name from URL path
func extractMethodName(path string) (string, error) {
	segments := strings.Split(path, "/")
	if len(segments) == 0 {
		return "", errors.New("invalid URL format")
	}
	return segments[len(segments)-1], nil
}

// Read JSON from request body into a given struct
func parseRequestBody(r *http.Request, target interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(target)
}

// Send JSON response to the client
func sendJSONResponse(w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func handleMethodCall(serviceVal reflect.Value, method reflect.Method, ctx context.Context, reqValue reflect.Value) ([]reflect.Value, error) {
	// Используем переданное значение напрямую вместо попытки развернуть
	return method.Func.Call([]reflect.Value{
		serviceVal,
		reflect.ValueOf(ctx),
		reqValue,
	}), nil
}

// Обновленная функция MakeHandler
func MakeHandler(service interface{}) http.Handler {
	serviceVal := reflect.ValueOf(service)
	serviceType := serviceVal.Type()

	methods := map[string]int{}
	for i := 0; i < serviceType.NumMethod(); i++ {
		method := serviceType.Method(i)
		methods[method.Name] = i
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		methodName, err := extractMethodName(r.URL.Path)
		if err != nil {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		methodIdx, exists := methods[methodName]
		if !exists {
			http.Error(w, "Method not found", http.StatusNotFound)
			return
		}

		method := serviceType.Method(methodIdx)

		reqVal := reflect.New(method.Type.In(2).Elem())
		if err := parseRequestBody(r, reqVal.Interface()); err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		results, err := handleMethodCall(serviceVal, method, context.Background(), reqVal)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if errVal := results[1].Interface(); errVal != nil {
			http.Error(w, errVal.(error).Error(), http.StatusInternalServerError)
			return
		}

		response := results[0].Interface()
		if err := sendJSONResponse(w, response); err != nil {
			http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		}
	})
}

// Call sends a JSON-RPC request and processes the response
func Call(ctx context.Context, endpoint, method string, req, rsp interface{}) error {
	url := fmt.Sprintf("%s/%s", endpoint, method)

	requestBytes, err := json.Marshal(req)
	if err != nil {
		return errors.New("failed to serialize request")
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBytes))
	if err != nil {
		return errors.New("failed to create HTTP request")
	}

	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return errors.New("request failed: " + err.Error())
	}

	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResp.Body)
		return errors.New("RPC error: " + string(body))
	}

	return json.NewDecoder(httpResp.Body).Decode(rsp)
}
