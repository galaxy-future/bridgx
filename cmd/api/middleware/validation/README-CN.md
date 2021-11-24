[[English](README.md)|中文]
# 校验器
## 自定义校验器
1. tag 名称
> 在结构体中添加 tag 并注册对应名称的校验器后可以进行对应的校验

example: tag `lt3`
```go
type Expmaple struct {
    Case string `validate:"lt3"`
}
```
2. 检验方法

example: lt3
```go
func lt3(fl validator.FieldLevel) bool {
    field := fl.Field().String()
    return len(field) < 3
}
```
3. 翻译方法
> 翻译方法主要是在校验失败后返回指定语言的错误信息(例如中文需要配合方法`Translate2Chinese`)

example:
```go
func translateLT3Err(ut ut.Translator, fe validator.FieldError) string {
    return "字符长度超过 3"
}
```
4. 注册校验器

在下列函数中添加 validation 即可.
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
5. 直接使用

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

6. gin 中的使用
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