import gleam/float
import gleam/io
import gleam/list
import nn/perceptron
import utils

pub fn run() {
  let divider = 3
  let bit_size = 4

  let train_data =
    [1, 3, 5, 6, 12, 14, 11, 23, 10, 22, 43, 32, 54, 21, 18, 27]
    |> list.map(fn(num) {
      let l = utils.conv_num(num, bit_size, divider)
      #(l, case num % 3 {
        0 -> 1.0
        _ -> 0.0
      })
    })
  // io.debug(train_data)

  let perc =
    list.fold(
      train_data,
      perceptron.new(bit_size, 1.0, utils.sig_act),
      fn(perc, pair) {
        let #(input, target) = pair
        perceptron.train(perc, input, target)
      },
    )

  io.println(
    "Weights "
    <> list.fold(perc.weights, "", fn(acc, w) {
      acc <> float.to_string(w) <> " "
    }),
  )

  let test_data =
    [8, 9, 6, 2]
    |> list.map(fn(num) {
      let l = utils.conv_num(num, bit_size, divider)
      #(l, case num % divider {
        0 -> 1.0
        _ -> 0.0
      })
    })
  // io.debug(test_data)

  list.each(test_data, fn(pair) {
    let #(input, target) = pair
    let res = perceptron.predict(perc, input)
    io.println(
      "Input [ "
      <> list.fold(input, "", fn(acc, part) {
        acc <> float.to_string(part) <> " "
      })
      <> "], Result "
      <> float.to_string(res)
      <> " , Target "
      <> float.to_string(target),
    )
  })
}
