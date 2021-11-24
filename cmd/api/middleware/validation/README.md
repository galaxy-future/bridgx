[English|[中文](README-CN.md)]
# Validation
## customer validator
1. Tag name
> After adding a tag to the structure and registering the validator
> of the corresponding name, the corresponding validation can be performed
example: tag `lt3`
```go
type Expmaple struct {
    Case string `validate:"lt3"`
}
```
2. Validation func

example: lt3
```go
func lt3(fl validator.FieldLevel) bool {
    field := fl.Field().String()
    return len(field) < 3
}

```
3. Translation func
> The translation method is mainly to return the error message of the specified
> language after the validation fails.(For example, Chinese needs to cooperate 
> with the method`Translate2Chinese`)

example:
```go
func translateLT3Err(ut ut.Translator, fe validator.FieldError) string {
    return "字符长度超过 3"
}
```
4. Register validator

Add your custom validation here.
```go
func RegisterCustomValidators() {
    appendMultiTagValidation(
        // Add your custom Validation here.
        Validation{
            validateFunc:     validateCharacterTypeGT3,
            translateFunc:    translateCharacterErr,
            translateRegFunc: defaultTranslateRegFunc,
            tag:              CharTypeGT3,
        },
        Validation{
            validateFunc:     lt3,
            translateFunc:    translateCharacterErr,
            translateRegFunc: translateLT3Err,
            tag:              "lt3",
        },
    )
}
```
5. Use directly

example:
```go
type ValidateCase struct {
    Name string `validate:"lt3"`
}

func main() {
    v := validator.New()
    RegisterCustomValidators()
    err := RegisterValidators(v)
    c := ValidateCase{Name: "1234"}
    err := v.Struct(c)
    fmt.Println(Translate2Chinese(err))
}
```
> 字符长度超过 3

6. Validation in gin 
```go
type validateCase struct {
    Name string `form:"name" binding:"required,lt3" json:"name"`
}

func main() {
    middleware.Init()
    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        err := validation.RegisterValidators(v)
        if err != nil {
        ...
        }
    }
}

// GET
func Getxx(ctx *gin.Context) {
    c := validateCase{}
    err := ctx.ShouldBindWith(&c, binding.Query)
    if err == nil {
    //"validate failed"
    }
}

// POST
func PostXXX(ctx *gin.Context) {
    c := validateCase{}
    err := ctx.ShouldBindJSON(&c)
    if err == nil {
    //"validate failed"
    }
}
```