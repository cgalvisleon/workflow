// Hello World
console.log("Hello World Step 0");
console.debug("Hello World Step 0");
console.info("Hello World Step 0");
console.error("Hello World Step 0");
console.log(
  "toString:",
  toString({
    ctx: ctx,
    ctxs: ctxs,
    pinnedData: pinnedData,
  })
);

let i = cache.incr("test", 10);
i = cache.incr("test", 10);
instance.SetPinnedData("test", "test2");

try {
  // model = model("octopus", "users");
  // console.log(model.ToJson().ToString());
} catch (error) {
  console.error(error);
}

result = {
  test: "test",
  test2: "test2",
  i: i,
};
