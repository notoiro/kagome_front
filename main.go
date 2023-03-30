package main

import (
	"strings"
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
	Tokens string `json:"tokens"`
}

func igOK(s string, _ bool) string {
	return s
}

func tokenize(str string) string {
	if len(str) == 0 {
		return ""
	}
	t, err := tokenizer.New(ipaneologd.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return ""
	}
	var ret []interface{}
	tokens := t.Tokenize(str)
	for _, v := range tokens {
		ret = append(ret, map[string]interface{}{
			"word_id":       v.ID,
			"word_type":     v.Class.String(),
			"word_position": v.Start,
			"surface_form":  v.Surface,
			"pos":           strings.Join(v.POS(), ","),
			"base_form":     igOK(v.BaseForm()),
			"reading":       igOK(v.Reading()),
			"pronunciation": igOK(v.Pronunciation()),
		})
	}
  var r, _ = json.Marshal(ret)
	return string(r)
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
		http.Error(w, fmt.Sprintf("{\"error\":\"%v\"}", err), http.StatusInternalServerError)
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

