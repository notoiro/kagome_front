package main

import (
  "fmt"
  "encoding/json"
  "net/http"
  "log"
  "flag"
  "strconv"

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

func start_server(port string) {
  mux := http.NewServeMux()
  mux.Handle("/tokenize", &TokenizeHandler{})

  srv := http.Server{
    Addr: ":" + port,
    Handler: mux,
  }

  log.Println(port);
  log.Println("localhost:" + port)
  srv.ListenAndServe()
}

func main() {
  var portFlag = flag.Int("port", 2971, "server port")
  flag.Parse()

  var port = strconv.Itoa(*portFlag)
  start_server(port)
  println("Start...")
}

