// @generated by protoc-gen-es v1.3.1 with parameter "target=ts"
// @generated from file rill/runtime/v1/catalog.proto (package rill.runtime.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3, Struct } from "@bufbuild/protobuf";
import { StructType } from "./schema_pb.js";
import { TimeGrain } from "./time_grain_pb.js";

/**
 * ObjectType represents the different kinds of catalog objects
 *
 * @generated from enum rill.runtime.v1.ObjectType
 */
export enum ObjectType {
  /**
   * @generated from enum value: OBJECT_TYPE_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * @generated from enum value: OBJECT_TYPE_TABLE = 1;
   */
  TABLE = 1,

  /**
   * @generated from enum value: OBJECT_TYPE_SOURCE = 2;
   */
  SOURCE = 2,

  /**
   * @generated from enum value: OBJECT_TYPE_MODEL = 3;
   */
  MODEL = 3,

  /**
   * @generated from enum value: OBJECT_TYPE_METRICS_VIEW = 4;
   */
  METRICS_VIEW = 4,
}
// Retrieve enum metadata with: proto3.getEnumType(ObjectType)
proto3.util.setEnumType(ObjectType, "rill.runtime.v1.ObjectType", [
  { no: 0, name: "OBJECT_TYPE_UNSPECIFIED" },
  { no: 1, name: "OBJECT_TYPE_TABLE" },
  { no: 2, name: "OBJECT_TYPE_SOURCE" },
  { no: 3, name: "OBJECT_TYPE_MODEL" },
  { no: 4, name: "OBJECT_TYPE_METRICS_VIEW" },
]);

/**
 * Table represents a table in the OLAP database. These include pre-existing tables discovered by periodically
 * scanning the database's information schema when the instance is created with exposed=true. Pre-existing tables
 * have managed = false.
 *
 * @generated from message rill.runtime.v1.Table
 */
export class Table extends Message<Table> {
  /**
   * Table name
   *
   * @generated from field: string name = 1;
   */
  name = "";

  /**
   * Table schema
   *
   * @generated from field: rill.runtime.v1.StructType schema = 2;
   */
  schema?: StructType;

  /**
   * Managed is true if the table was created through a runtime migration, false if it was discovered in by
   * scanning the database's information schema.
   *
   * @generated from field: bool managed = 3;
   */
  managed = false;

  constructor(data?: PartialMessage<Table>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.runtime.v1.Table";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "schema", kind: "message", T: StructType },
    { no: 3, name: "managed", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Table {
    return new Table().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Table {
    return new Table().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Table {
    return new Table().fromJsonString(jsonString, options);
  }

  static equals(a: Table | PlainMessage<Table> | undefined, b: Table | PlainMessage<Table> | undefined): boolean {
    return proto3.util.equals(Table, a, b);
  }
}

/**
 * Source is the internal representation of a source definition
 *
 * @generated from message rill.runtime.v1.Source
 */
export class Source extends Message<Source> {
  /**
   * Name of the source
   *
   * @generated from field: string name = 1;
   */
  name = "";

  /**
   * Connector used by the source
   *
   * @generated from field: string connector = 2;
   */
  connector = "";

  /**
   * Connector properties assigned in the source
   *
   * @generated from field: google.protobuf.Struct properties = 3;
   */
  properties?: Struct;

  /**
   * Detected schema of the source
   *
   * @generated from field: rill.runtime.v1.StructType schema = 5;
   */
  schema?: StructType;

  /**
   * timeout for source ingestion in seconds
   *
   * @generated from field: int32 timeout_seconds = 7;
   */
  timeoutSeconds = 0;

  constructor(data?: PartialMessage<Source>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.runtime.v1.Source";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "connector", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "properties", kind: "message", T: Struct },
    { no: 5, name: "schema", kind: "message", T: StructType },
    { no: 7, name: "timeout_seconds", kind: "scalar", T: 5 /* ScalarType.INT32 */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Source {
    return new Source().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Source {
    return new Source().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Source {
    return new Source().fromJsonString(jsonString, options);
  }

  static equals(a: Source | PlainMessage<Source> | undefined, b: Source | PlainMessage<Source> | undefined): boolean {
    return proto3.util.equals(Source, a, b);
  }
}

/**
 * Model is the internal representation of a model definition
 *
 * @generated from message rill.runtime.v1.Model
 */
export class Model extends Message<Model> {
  /**
   * Name of the model
   *
   * @generated from field: string name = 1;
   */
  name = "";

  /**
   * SQL is a SELECT statement representing the model
   *
   * @generated from field: string sql = 2;
   */
  sql = "";

  /**
   * Dialect of the SQL statement
   *
   * @generated from field: rill.runtime.v1.Model.Dialect dialect = 3;
   */
  dialect = Model_Dialect.UNSPECIFIED;

  /**
   * Detected schema of the model
   *
   * @generated from field: rill.runtime.v1.StructType schema = 4;
   */
  schema?: StructType;

  /**
   * To materialize model or not
   *
   * @generated from field: bool materialize = 5;
   */
  materialize = false;

  constructor(data?: PartialMessage<Model>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.runtime.v1.Model";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "sql", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "dialect", kind: "enum", T: proto3.getEnumType(Model_Dialect) },
    { no: 4, name: "schema", kind: "message", T: StructType },
    { no: 5, name: "materialize", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Model {
    return new Model().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Model {
    return new Model().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Model {
    return new Model().fromJsonString(jsonString, options);
  }

  static equals(a: Model | PlainMessage<Model> | undefined, b: Model | PlainMessage<Model> | undefined): boolean {
    return proto3.util.equals(Model, a, b);
  }
}

/**
 * Dialects supported for models
 *
 * @generated from enum rill.runtime.v1.Model.Dialect
 */
export enum Model_Dialect {
  /**
   * @generated from enum value: DIALECT_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * @generated from enum value: DIALECT_DUCKDB = 1;
   */
  DUCKDB = 1,
}
// Retrieve enum metadata with: proto3.getEnumType(Model_Dialect)
proto3.util.setEnumType(Model_Dialect, "rill.runtime.v1.Model.Dialect", [
  { no: 0, name: "DIALECT_UNSPECIFIED" },
  { no: 1, name: "DIALECT_DUCKDB" },
]);

/**
 * Metrics view is the internal representation of a metrics view definition
 *
 * @generated from message rill.runtime.v1.MetricsView
 */
export class MetricsView extends Message<MetricsView> {
  /**
   * Name of the metrics view
   *
   * @generated from field: string name = 1;
   */
  name = "";

  /**
   * Name of the source or model that the metrics view is based on
   *
   * @generated from field: string model = 2;
   */
  model = "";

  /**
   * Name of the primary time dimension, used for rendering time series
   *
   * @generated from field: string time_dimension = 3;
   */
  timeDimension = "";

  /**
   * Dimensions in the metrics view
   *
   * @generated from field: repeated rill.runtime.v1.MetricsView.Dimension dimensions = 5;
   */
  dimensions: MetricsView_Dimension[] = [];

  /**
   * Measures in the metrics view
   *
   * @generated from field: repeated rill.runtime.v1.MetricsView.Measure measures = 6;
   */
  measures: MetricsView_Measure[] = [];

  /**
   * User friendly label for the dashboard
   *
   * @generated from field: string label = 7;
   */
  label = "";

  /**
   * Brief description of the dashboard
   *
   * @generated from field: string description = 8;
   */
  description = "";

  /**
   * Smallest time grain to show in the dashboard
   *
   * @generated from field: rill.runtime.v1.TimeGrain smallest_time_grain = 9;
   */
  smallestTimeGrain = TimeGrain.UNSPECIFIED;

  /**
   * Default time range for the dashboard. It should be a valid ISO 8601 duration string.
   *
   * @generated from field: string default_time_range = 10;
   */
  defaultTimeRange = "";

  /**
   * Available time zones list preferred time zones using IANA location identifiers.
   *
   * @generated from field: repeated string available_time_zones = 11;
   */
  availableTimeZones: string[] = [];

  /**
   * Security for the dashboard
   *
   * @generated from field: rill.runtime.v1.MetricsView.Security security = 12;
   */
  security?: MetricsView_Security;

  /**
   * @generated from field: uint32 first_day_of_week = 13;
   */
  firstDayOfWeek = 0;

  /**
   * @generated from field: uint32 first_month_of_year = 14;
   */
  firstMonthOfYear = 0;

  constructor(data?: PartialMessage<MetricsView>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.runtime.v1.MetricsView";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "model", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "time_dimension", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 5, name: "dimensions", kind: "message", T: MetricsView_Dimension, repeated: true },
    { no: 6, name: "measures", kind: "message", T: MetricsView_Measure, repeated: true },
    { no: 7, name: "label", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 8, name: "description", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 9, name: "smallest_time_grain", kind: "enum", T: proto3.getEnumType(TimeGrain) },
    { no: 10, name: "default_time_range", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 11, name: "available_time_zones", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
    { no: 12, name: "security", kind: "message", T: MetricsView_Security },
    { no: 13, name: "first_day_of_week", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 14, name: "first_month_of_year", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MetricsView {
    return new MetricsView().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MetricsView {
    return new MetricsView().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MetricsView {
    return new MetricsView().fromJsonString(jsonString, options);
  }

  static equals(a: MetricsView | PlainMessage<MetricsView> | undefined, b: MetricsView | PlainMessage<MetricsView> | undefined): boolean {
    return proto3.util.equals(MetricsView, a, b);
  }
}

/**
 * Dimensions are columns to filter and group by
 *
 * @generated from message rill.runtime.v1.MetricsView.Dimension
 */
export class MetricsView_Dimension extends Message<MetricsView_Dimension> {
  /**
   * @generated from field: string name = 1;
   */
  name = "";

  /**
   * @generated from field: string label = 2;
   */
  label = "";

  /**
   * @generated from field: string description = 3;
   */
  description = "";

  /**
   * @generated from field: string column = 4;
   */
  column = "";

  constructor(data?: PartialMessage<MetricsView_Dimension>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.runtime.v1.MetricsView.Dimension";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "label", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "description", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "column", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MetricsView_Dimension {
    return new MetricsView_Dimension().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MetricsView_Dimension {
    return new MetricsView_Dimension().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MetricsView_Dimension {
    return new MetricsView_Dimension().fromJsonString(jsonString, options);
  }

  static equals(a: MetricsView_Dimension | PlainMessage<MetricsView_Dimension> | undefined, b: MetricsView_Dimension | PlainMessage<MetricsView_Dimension> | undefined): boolean {
    return proto3.util.equals(MetricsView_Dimension, a, b);
  }
}

/**
 * Measures are aggregated computed values
 *
 * @generated from message rill.runtime.v1.MetricsView.Measure
 */
export class MetricsView_Measure extends Message<MetricsView_Measure> {
  /**
   * @generated from field: string name = 1;
   */
  name = "";

  /**
   * @generated from field: string label = 2;
   */
  label = "";

  /**
   * @generated from field: string expression = 3;
   */
  expression = "";

  /**
   * @generated from field: string description = 4;
   */
  description = "";

  /**
   * @generated from field: string format = 5;
   */
  format = "";

  /**
   * @generated from field: bool valid_percent_of_total = 6;
   */
  validPercentOfTotal = false;

  constructor(data?: PartialMessage<MetricsView_Measure>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.runtime.v1.MetricsView.Measure";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "label", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "expression", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "description", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 5, name: "format", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 6, name: "valid_percent_of_total", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MetricsView_Measure {
    return new MetricsView_Measure().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MetricsView_Measure {
    return new MetricsView_Measure().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MetricsView_Measure {
    return new MetricsView_Measure().fromJsonString(jsonString, options);
  }

  static equals(a: MetricsView_Measure | PlainMessage<MetricsView_Measure> | undefined, b: MetricsView_Measure | PlainMessage<MetricsView_Measure> | undefined): boolean {
    return proto3.util.equals(MetricsView_Measure, a, b);
  }
}

/**
 * Security for the metrics view
 *
 * @generated from message rill.runtime.v1.MetricsView.Security
 */
export class MetricsView_Security extends Message<MetricsView_Security> {
  /**
   * Dashboard level access condition
   *
   * @generated from field: string access = 1;
   */
  access = "";

  /**
   * row level access condition
   *
   * @generated from field: string row_filter = 2;
   */
  rowFilter = "";

  /**
   * either one of include or exclude will be specified
   *
   * @generated from field: repeated rill.runtime.v1.MetricsView.Security.FieldCondition include = 3;
   */
  include: MetricsView_Security_FieldCondition[] = [];

  /**
   * @generated from field: repeated rill.runtime.v1.MetricsView.Security.FieldCondition exclude = 4;
   */
  exclude: MetricsView_Security_FieldCondition[] = [];

  constructor(data?: PartialMessage<MetricsView_Security>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.runtime.v1.MetricsView.Security";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "access", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "row_filter", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "include", kind: "message", T: MetricsView_Security_FieldCondition, repeated: true },
    { no: 4, name: "exclude", kind: "message", T: MetricsView_Security_FieldCondition, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MetricsView_Security {
    return new MetricsView_Security().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MetricsView_Security {
    return new MetricsView_Security().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MetricsView_Security {
    return new MetricsView_Security().fromJsonString(jsonString, options);
  }

  static equals(a: MetricsView_Security | PlainMessage<MetricsView_Security> | undefined, b: MetricsView_Security | PlainMessage<MetricsView_Security> | undefined): boolean {
    return proto3.util.equals(MetricsView_Security, a, b);
  }
}

/**
 * Dimension/measure access condition
 *
 * @generated from message rill.runtime.v1.MetricsView.Security.FieldCondition
 */
export class MetricsView_Security_FieldCondition extends Message<MetricsView_Security_FieldCondition> {
  /**
   * @generated from field: string condition = 1;
   */
  condition = "";

  /**
   * @generated from field: repeated string names = 2;
   */
  names: string[] = [];

  constructor(data?: PartialMessage<MetricsView_Security_FieldCondition>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "rill.runtime.v1.MetricsView.Security.FieldCondition";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "condition", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "names", kind: "scalar", T: 9 /* ScalarType.STRING */, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MetricsView_Security_FieldCondition {
    return new MetricsView_Security_FieldCondition().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MetricsView_Security_FieldCondition {
    return new MetricsView_Security_FieldCondition().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MetricsView_Security_FieldCondition {
    return new MetricsView_Security_FieldCondition().fromJsonString(jsonString, options);
  }

  static equals(a: MetricsView_Security_FieldCondition | PlainMessage<MetricsView_Security_FieldCondition> | undefined, b: MetricsView_Security_FieldCondition | PlainMessage<MetricsView_Security_FieldCondition> | undefined): boolean {
    return proto3.util.equals(MetricsView_Security_FieldCondition, a, b);
  }
}

