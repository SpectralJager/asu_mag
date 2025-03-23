import gleam/float
import gleam/list

pub type Perceptron {
  Perceptron(weights: List(Float), rate: Float, act_fn: fn(Float) -> Float)
}

pub fn new_perceptron(
  input_size: Int,
  learning_rate: Float,
  activation_fn: fn(Float) -> Float,
) {
  Perceptron(
    list.repeat(0.5, input_size),
    // |> list.map(fn(_) { float.random() }),
    learning_rate,
    activation_fn,
  )
}

pub fn predict(perceptron: Perceptron, input: List(Float)) -> Float {
  let net = sum(perceptron.weights, input)
  perceptron.act_fn(net)
}

pub fn train(
  perceptron: Perceptron,
  input: List(Float),
  target: Float,
) -> Perceptron {
  let res = predict(perceptron, input)
  let error = target -. res
  let weights =
    list.map2(perceptron.weights, input, fn(w, in) {
      w +. perceptron.rate *. error *. in
    })
  Perceptron(..perceptron, weights: weights)
}

fn sum(weights: List(Float), inputs: List(Float)) -> Float {
  list.map2(weights, inputs, fn(w, in) { #(w, in) })
  |> list.fold(0.0, fn(acc, pair) {
    let #(w, in) = pair
    acc +. w *. in
  })
}
