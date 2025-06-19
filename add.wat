(module
  (func (export "add") (param $a i32) (param $b i32) (result i32)
    local.get $a  ;; Push the first parameter ($a) onto the stack.
    local.get $b  ;; Push the second parameter ($b) onto the stack.
    i32.add      ;; Pop the top two values, add them, and push the result onto the stack.
  )
)