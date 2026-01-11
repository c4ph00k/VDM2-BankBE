//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class Transfer {
  /// Returns a new [Transfer] instance.
  Transfer({
    required this.id,
    required this.fromAccount,
    required this.toAccount,
    required this.amount,
    required this.status,
    required this.initiatedAt,
    required this.completedAt,
  });

  int id;

  String fromAccount;

  String toAccount;

  /// Decimal encoded as string (shopspring/decimal)
  String amount;

  TransferStatusEnum status;

  DateTime initiatedAt;

  DateTime? completedAt;

  @override
  bool operator ==(Object other) => identical(this, other) || other is Transfer &&
    other.id == id &&
    other.fromAccount == fromAccount &&
    other.toAccount == toAccount &&
    other.amount == amount &&
    other.status == status &&
    other.initiatedAt == initiatedAt &&
    other.completedAt == completedAt;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (id.hashCode) +
    (fromAccount.hashCode) +
    (toAccount.hashCode) +
    (amount.hashCode) +
    (status.hashCode) +
    (initiatedAt.hashCode) +
    (completedAt == null ? 0 : completedAt!.hashCode);

  @override
  String toString() => 'Transfer[id=$id, fromAccount=$fromAccount, toAccount=$toAccount, amount=$amount, status=$status, initiatedAt=$initiatedAt, completedAt=$completedAt]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'id'] = this.id;
      json[r'from_account'] = this.fromAccount;
      json[r'to_account'] = this.toAccount;
      json[r'amount'] = this.amount;
      json[r'status'] = this.status;
      json[r'initiated_at'] = this.initiatedAt.toUtc().toIso8601String();
    if (this.completedAt != null) {
      json[r'completed_at'] = this.completedAt!.toUtc().toIso8601String();
    } else {
      json[r'completed_at'] = null;
    }
    return json;
  }

  /// Returns a new [Transfer] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static Transfer? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "Transfer[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "Transfer[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return Transfer(
        id: mapValueOfType<int>(json, r'id')!,
        fromAccount: mapValueOfType<String>(json, r'from_account')!,
        toAccount: mapValueOfType<String>(json, r'to_account')!,
        amount: mapValueOfType<String>(json, r'amount')!,
        status: TransferStatusEnum.fromJson(json[r'status'])!,
        initiatedAt: mapDateTime(json, r'initiated_at', r'')!,
        completedAt: mapDateTime(json, r'completed_at', r''),
      );
    }
    return null;
  }

  static List<Transfer> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <Transfer>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = Transfer.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, Transfer> mapFromJson(dynamic json) {
    final map = <String, Transfer>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = Transfer.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of Transfer-objects as value to a dart map
  static Map<String, List<Transfer>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<Transfer>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = Transfer.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'id',
    'from_account',
    'to_account',
    'amount',
    'status',
    'initiated_at',
    'completed_at',
  };
}


class TransferStatusEnum {
  /// Instantiate a new enum with the provided [value].
  const TransferStatusEnum._(this.value);

  /// The underlying value of this enum member.
  final String value;

  @override
  String toString() => value;

  String toJson() => value;

  static const pending = TransferStatusEnum._(r'pending');
  static const completed = TransferStatusEnum._(r'completed');
  static const failed = TransferStatusEnum._(r'failed');

  /// List of all possible values in this [enum][TransferStatusEnum].
  static const values = <TransferStatusEnum>[
    pending,
    completed,
    failed,
  ];

  static TransferStatusEnum? fromJson(dynamic value) => TransferStatusEnumTypeTransformer().decode(value);

  static List<TransferStatusEnum> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <TransferStatusEnum>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = TransferStatusEnum.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }
}

/// Transformation class that can [encode] an instance of [TransferStatusEnum] to String,
/// and [decode] dynamic data back to [TransferStatusEnum].
class TransferStatusEnumTypeTransformer {
  factory TransferStatusEnumTypeTransformer() => _instance ??= const TransferStatusEnumTypeTransformer._();

  const TransferStatusEnumTypeTransformer._();

  String encode(TransferStatusEnum data) => data.value;

  /// Decodes a [dynamic value][data] to a TransferStatusEnum.
  ///
  /// If [allowNull] is true and the [dynamic value][data] cannot be decoded successfully,
  /// then null is returned. However, if [allowNull] is false and the [dynamic value][data]
  /// cannot be decoded successfully, then an [UnimplementedError] is thrown.
  ///
  /// The [allowNull] is very handy when an API changes and a new enum value is added or removed,
  /// and users are still using an old app with the old code.
  TransferStatusEnum? decode(dynamic data, {bool allowNull = true}) {
    if (data != null) {
      switch (data) {
        case r'pending': return TransferStatusEnum.pending;
        case r'completed': return TransferStatusEnum.completed;
        case r'failed': return TransferStatusEnum.failed;
        default:
          if (!allowNull) {
            throw ArgumentError('Unknown enum value to decode: $data');
          }
      }
    }
    return null;
  }

  /// Singleton [TransferStatusEnumTypeTransformer] instance.
  static TransferStatusEnumTypeTransformer? _instance;
}


