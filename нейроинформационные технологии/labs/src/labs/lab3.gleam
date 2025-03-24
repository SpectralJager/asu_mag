import gleam/io
import nn/network

pub fn run() {
  let net = network.new([2, 1])
  io.debug(net)
  let res = network.predict(net, [1.0, 1.0])
  io.debug(res)
}
