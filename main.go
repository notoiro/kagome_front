package main

import (
  "fmt"
  "encoding/json"
  "net/http"
  "log"

  "github.com/ikawaha/kagome-dict-ipa-neologd"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

type TokenizeHandler struct{}
type TokenizerRequestBody struct {
	Input string `json:"text"`
}
type TokenizerResponseBody struct {
	Tokens []tokenizer.TokenData `json:"tokens"`
}

func igOK(s string, _ bool) string {
	return s
}

func tokenize(str string) []tokenizer.TokenData {
	if len(str) == 0 {
		return nil
	}
	t, err := tokenizer.New(ipaneologd.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return nil
	}
	tokens := t.Tokenize(str)
  var tokenData []tokenizer.TokenData
  for _, v := range tokens{
    tokenData = append(tokenData, tokenizer.NewTokenData(v))
  }
  return tokenData
}

func (h *TokenizeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "application/json")
	var req TokenizerRequestBody
  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("{\"error\":\"%v\"}", err), http.StatusBadRequest)
		return
	}

  if req.Input == "" {
		w.Write([]byte(`{"tokens":[]}`))
		return
	}

  tokens := tokenize(req.Input)

	resp, err := json.Marshal(TokenizerResponseBody{
		Tokens: tokens,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("{\"status\":false,\"error\":\"%v\"}", err), http.StatusInternalServerError)
		return
	}
  w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func start_server() {
  mux := http.NewServeMux()
  mux.Handle("/tokenize", &TokenizeHandler{})

  srv := http.Server{
    Addr: ":2971",
    Handler: mux,
  }

  log.Println("localhost:2971")
  srv.ListenAndServe()
}

func main() {
  start_server()
  println("Start...")
}

