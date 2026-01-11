//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class CreateMovementRequest {
  /// Returns a new [CreateMovementRequest] instance.
  CreateMovementRequest({
    required this.amount,
    required this.type,
    this.description,
  });

  /// Decimal encoded as string (shopspring/decimal)
  String amount;

  CreateMovementRequestTypeEnum type;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? description;

  @override
  bool operator ==(Object other) => identical(this, other) || other is CreateMovementRequest &&
    other.amount == amount &&
    other.type == type &&
    other.description == description;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (amount.hashCode) +
    (type.hashCode) +
    (description == null ? 0 : description!.hashCode);

  @override
  String toString() => 'CreateMovementRequest[amount=$amount, type=$type, description=$description]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'amount'] = this.amount;
      json[r'type'] = this.type;
    if (this.description != null) {
      json[r'description'] = this.description;
    } else {
      json[r'description'] = null;
    }
    return json;
  }

  /// Returns a new [CreateMovementRequest] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static CreateMovementRequest? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "CreateMovementRequest[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "CreateMovementRequest[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return CreateMovementRequest(
        amount: mapValueOfType<String>(json, r'amount')!,
        type: CreateMovementRequestTypeEnum.fromJson(json[r'type'])!,
        description: mapValueOfType<String>(json, r'description'),
      );
    }
    return null;
  }

  static List<CreateMovementRequest> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <CreateMovementRequest>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = CreateMovementRequest.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, CreateMovementRequest> mapFromJson(dynamic json) {
    final map = <String, CreateMovementRequest>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = CreateMovementRequest.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of CreateMovementRequest-objects as value to a dart map
  static Map<String, List<CreateMovementRequest>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<CreateMovementRequest>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = CreateMovementRequest.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'amount',
    'type',
  };
}


class CreateMovementRequestTypeEnum {
  /// Instantiate a new enum with the provided [value].
  const CreateMovementRequestTypeEnum._(this.value);

  /// The underlying value of this enum member.
  final String value;

  @override
  String toString() => value;

  String toJson() => value;

  static const credit = CreateMovementRequestTypeEnum._(r'credit');
  static const debit = CreateMovementRequestTypeEnum._(r'debit');

  /// List of all possible values in this [enum][CreateMovementRequestTypeEnum].
  static const values = <CreateMovementRequestTypeEnum>[
    credit,
    debit,
  ];

  static CreateMovementRequestTypeEnum? fromJson(dynamic value) => CreateMovementRequestTypeEnumTypeTransformer().decode(value);

  static List<CreateMovementRequestTypeEnum> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <CreateMovementRequestTypeEnum>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = CreateMovementRequestTypeEnum.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }
}

/// Transformation class that can [encode] an instance of [CreateMovementRequestTypeEnum] to String,
/// and [decode] dynamic data back to [CreateMovementRequestTypeEnum].
class CreateMovementRequestTypeEnumTypeTransformer {
  factory CreateMovementRequestTypeEnumTypeTransformer() => _instance ??= const CreateMovementRequestTypeEnumTypeTransformer._();

  const CreateMovementRequestTypeEnumTypeTransformer._();

  String encode(CreateMovementRequestTypeEnum data) => data.value;

  /// Decodes a [dynamic value][data] to a CreateMovementRequestTypeEnum.
  ///
  /// If [allowNull] is true and the [dynamic value][data] cannot be decoded successfully,
  /// then null is returned. However, if [allowNull] is false and the [dynamic value][data]
  /// cannot be decoded successfully, then an [UnimplementedError] is thrown.
  ///
  /// The [allowNull] is very handy when an API changes and a new enum value is added or removed,
  /// and users are still using an old app with the old code.
  CreateMovementRequestTypeEnum? decode(dynamic data, {bool allowNull = true}) {
    if (data != null) {
      switch (data) {
        case r'credit': return CreateMovementRequestTypeEnum.credit;
        case r'debit': return CreateMovementRequestTypeEnum.debit;
        default:
          if (!allowNull) {
            throw ArgumentError('Unknown enum value to decode: $data');
          }
      }
    }
    return null;
  }

  /// Singleton [CreateMovementRequestTypeEnumTypeTransformer] instance.
  static CreateMovementRequestTypeEnumTypeTransformer? _instance;
}


