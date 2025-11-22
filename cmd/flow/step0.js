// Hello World
console.log("Hello World");
console.debug("Hello World");
console.info("Hello World");
console.error("Hello World");
console.log({
  ctx: ctx,
  ctxs: ctxs,
  pinnedData: pinnedData,
});

let i = cache.incr("test", 10);
i = cache.incr("test", 10);
console.log(i);
instance.SetPinnedData("test", "test2");

result = {
  test: "test",
  test2: "test2",
};
