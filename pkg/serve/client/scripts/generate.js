'use strict';
var fs = require('fs');
var CodeGen = require("swagger-typescript-codegen").CodeGen;

var swaggerFile = '../api/swagger.json';
var swagger = JSON.parse(fs.readFileSync(swaggerFile, "UTF-8"));
var tsCode = CodeGen.getTypescriptCode({
  className: "Api",
  swagger: swagger,
});
fs.mkdirSync('src/generated', { recursive: true });
fs.writeFileSync('src/generated/api.ts', tsCode, { encoding: "utf8", mode: 0o644});
