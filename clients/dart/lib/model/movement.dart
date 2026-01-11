//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class Movement {
  /// Returns a new [Movement] instance.
  Movement({
    required this.id,
    required this.accountId,
    required this.amount,
    required this.type,
    required this.description,
    required this.occurredAt,
  });

  int id;

  String accountId;

  /// Decimal encoded as string (shopspring/decimal)
  String amount;

  MovementTypeEnum type;

  String description;

  DateTime occurredAt;

  @override
  bool operator ==(Object other) => identical(this, other) || other is Movement &&
    other.id == id &&
    other.accountId == accountId &&
    other.amount == amount &&
    other.type == type &&
    other.description == description &&
    other.occurredAt == occurredAt;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (id.hashCode) +
    (accountId.hashCode) +
    (amount.hashCode) +
    (type.hashCode) +
    (description.hashCode) +
    (occurredAt.hashCode);

  @override
  String toString() => 'Movement[id=$id, accountId=$accountId, amount=$amount, type=$type, description=$description, occurredAt=$occurredAt]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'id'] = this.id;
      json[r'account_id'] = this.accountId;
      json[r'amount'] = this.amount;
      json[r'type'] = this.type;
      json[r'description'] = this.description;
      json[r'occurred_at'] = this.occurredAt.toUtc().toIso8601String();
    return json;
  }

  /// Returns a new [Movement] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static Movement? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "Movement[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "Movement[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return Movement(
        id: mapValueOfType<int>(json, r'id')!,
        accountId: mapValueOfType<String>(json, r'account_id')!,
        amount: mapValueOfType<String>(json, r'amount')!,
        type: MovementTypeEnum.fromJson(json[r'type'])!,
        description: mapValueOfType<String>(json, r'description')!,
        occurredAt: mapDateTime(json, r'occurred_at', r'')!,
      );
    }
    return null;
  }

  static List<Movement> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <Movement>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = Movement.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, Movement> mapFromJson(dynamic json) {
    final map = <String, Movement>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = Movement.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of Movement-objects as value to a dart map
  static Map<String, List<Movement>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<Movement>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = Movement.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'id',
    'account_id',
    'amount',
    'type',
    'description',
    'occurred_at',
  };
}


class MovementTypeEnum {
  /// Instantiate a new enum with the provided [value].
  const MovementTypeEnum._(this.value);

  /// The underlying value of this enum member.
  final String value;

  @override
  String toString() => value;

  String toJson() => value;

  static const credit = MovementTypeEnum._(r'credit');
  static const debit = MovementTypeEnum._(r'debit');

  /// List of all possible values in this [enum][MovementTypeEnum].
  static const values = <MovementTypeEnum>[
    credit,
    debit,
  ];

  static MovementTypeEnum? fromJson(dynamic value) => MovementTypeEnumTypeTransformer().decode(value);

  static List<MovementTypeEnum> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <MovementTypeEnum>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = MovementTypeEnum.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }
}

/// Transformation class that can [encode] an instance of [MovementTypeEnum] to String,
/// and [decode] dynamic data back to [MovementTypeEnum].
class MovementTypeEnumTypeTransformer {
  factory MovementTypeEnumTypeTransformer() => _instance ??= const MovementTypeEnumTypeTransformer._();

  const MovementTypeEnumTypeTransformer._();

  String encode(MovementTypeEnum data) => data.value;

  /// Decodes a [dynamic value][data] to a MovementTypeEnum.
  ///
  /// If [allowNull] is true and the [dynamic value][data] cannot be decoded successfully,
  /// then null is returned. However, if [allowNull] is false and the [dynamic value][data]
  /// cannot be decoded successfully, then an [UnimplementedError] is thrown.
  ///
  /// The [allowNull] is very handy when an API changes and a new enum value is added or removed,
  /// and users are still using an old app with the old code.
  MovementTypeEnum? decode(dynamic data, {bool allowNull = true}) {
    if (data != null) {
      switch (data) {
        case r'credit': return MovementTypeEnum.credit;
        case r'debit': return MovementTypeEnum.debit;
        default:
          if (!allowNull) {
            throw ArgumentError('Unknown enum value to decode: $data');
          }
      }
    }
    return null;
  }

  /// Singleton [MovementTypeEnumTypeTransformer] instance.
  static MovementTypeEnumTypeTransformer? _instance;
}


