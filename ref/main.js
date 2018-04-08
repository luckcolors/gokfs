'use strict';

const utils = require('./utils');

console.log(utils.isValidKey(""));
console.log(utils.hashKey(""));
console.log(utils.coerceKey(""));
console.log(utils.createItemKeyFromIndex("", 0));
console.log(utils.createSbucketNameFromIndex(0));
console.log(utils.createReferenceId(0));
console.log(utils.fileDoesExist("./main.js"));
console.log(utils.toHumanReadableSize(100000));
console.log(utils.coerceTablePath(""));
console.log(utils.isNotFoundError(0));
console.log("\n");
console.log("\n");
console.log("\n");
console.log(utils.isValidKey("A"))
console.log(utils.coerceKey("A"));
console.log(utils.isValidKey("ddadef707ba62c166051b9e3cd0294c27515f2bc"));
console.log(utils.hashKey("A"));
console.log("54e15c41a9dea039e790a69ce64ca4077a980fcc")
console.log(utils.coerceKey("A"));
console.log(utils.createItemKeyFromIndex("ddadef707ba62c166051b9e3cd0294c27515f2bc", 2213));
console.log(utils.createSbucketNameFromIndex(2213));
console.log(utils.createReferenceId());
console.log(utils.fileDoesExist("./main.js"));
console.log(utils.toHumanReadableSize(100000));
console.log(utils.coerceTablePath("a"));
console.log(utils.coerceTablePath("a.kfs"));
console.log(utils.isNotFoundError(0));
