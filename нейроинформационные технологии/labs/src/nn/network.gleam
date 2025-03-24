import gleam/io
import gleam/list
import nn/perceptron as p
import utils

pub type Network {
  Network(layers: List(List(p.Perceptron)))
}

pub fn new(layers_size: List(Int)) -> Network {
  Network(new_loop([], layers_size, 1))
}

fn new_loop(
  layers: List(List(p.Perceptron)),
  sizes: List(Int),
  prev_size: Int,
) -> List(List(p.Perceptron)) {
  case sizes {
    [] -> layers
    [size, ..next] -> {
      let layers =
        list.append(layers, [
          list.repeat(p.new(prev_size, 1.0, utils.sig_act), size),
        ])
      new_loop(layers, next, size)
    }
  }
}

pub fn predict(network: Network, input: List(Float)) -> List(Float) {
  predict_loop(input, network.layers, True)
}

fn predict_loop(
  input: List(Float),
  layers: List(List(p.Perceptron)),
  is_input_layer: Bool,
) -> List(Float) {
  case layers {
    [] -> input
    [layer, ..next] if is_input_layer -> {
      let res = list.map2(layer, input, fn(per, in) { p.predict(per, [in]) })
      predict_loop(res, next, False)
    }
    [layer, ..next] -> {
      let res = list.map(layer, fn(per) { p.predict(per, input) })
      predict_loop(res, next, False)
    }
  }
}

pub fn train(_input: List(Float), _train_fn: fn() -> Float) -> Network {
  todo
}

fn train_loop() {
  todo
}
