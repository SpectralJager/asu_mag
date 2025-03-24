import gleam/int
import gleam/list

pub fn conv_num(num: Int, size: Int, divider: Int) -> List(Float) {
  list.repeat(0.0, size)
  |> list.index_map(fn(_, i) {
    case i {
      i if i == 0 -> int.to_float(num % divider)
      i ->
        int.to_float({
          num / { int.product(list.repeat(divider, i)) } % divider
        })
    }
  })
  |> list.reverse
}

pub fn sig_act(net: Float) -> Float {
  case net {
    a if a >=. 0.5 -> 1.0
    _ -> 0.0
  }
}
