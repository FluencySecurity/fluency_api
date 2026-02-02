= Fluency API overview

* All APIs are HTTP POST => https://$domain/prefix/$function
**prefixs:  "api/auth" (user/role related)  or "api/ds" (other apis)
**The function call request format: fsb.CallRequest 
[source, golang]
type CallRequest struct {
   Function    string        `json:"function"`
   Kargs       json.RawMessage        `json:"kargs"`
}
** The response format: fsb.CallResponse
***if the verdict is "OK", the response field contains the result.
***if the verdict is "ERROR", the error field show the error information.
[source, golang]
type CallResponse struct {
   Verdict string `json:"verdict"`
   Response    json.RawMessage        `json:"response"`
   Error       string        `json:"error,omitempty"`
}
* For typical CRUD operations, the dao APIs (with a suffix "_dao") has the following format. 
** action values: get, add, delete, update, list, toggle
[source, golang] 
type GenericDaoRequestArgs[T any] struct {
  Id    string `json:"id,omitempty"`
  Entry *T      `json:"entry,omitempty"`
  Flag  bool   `json:"flag,omitempty"`  // optional for toggle action
}
type GenericDaoRequest[T any] struct {
	Action string                    `json:"action"`
	Args   *GenericDaoRequestArgs[T] `json:"args,omitempty"`
}

== authorization related APIs: (with prefix 'api/auth')

* userAdd
* userList
* getUser
* userDelete

==  
