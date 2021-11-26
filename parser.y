%{
package variable_json

type pair struct {
  key string
  val interface{}
}
%}

%union{
  obj map[string]interface{}
  list []interface{}
  pair pair
  val interface{}
}

%token LexError
%token <val> String Number Literal '{' '}' '[' ']' Variable

%type <obj> members
%type <pair> pair
%type <list> elements
%type <val> value


%start value

%%
members:
  {
    $$ = map[string]interface{}{}
  }
| pair
  {
    $$ = map[string]interface{}{
      $1.key: $1.val,
    }
  }
| members ',' pair
  {
    $1[$3.key] = $3.val
    $$ = $1
  }

pair: String ':' value
  {
    $$ = pair{key: $1.(string), val: $3}
  }

elements:
  {
    $$ = []interface{}{}
  }
| value
  {
    $$ = []interface{}{$1}
  }
| elements ',' value
  {
    $$ = append($1, $3)
  }

value:
  String
  {
     $$ = newString($1)
     maybeSetResult(yylex, $$)
  }
| Number
  {
    $$ = newNumber($1)
    maybeSetResult(yylex, $$)
  }
| Literal
  {
    $$ = newLiteral($1)
    maybeSetResult(yylex, $$)
  }
| Variable
  {
    $$ = newVariable($1)
    maybeSetResult(yylex, $$)
  }
| '{' members '}'
  {
    $$ = newObject($2)
    maybeSetResult(yylex, $$)
  }
| '[' elements ']'
  {
    $$ = newArray($2)
    maybeSetResult(yylex, $$)
  }