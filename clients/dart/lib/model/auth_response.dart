//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class AuthResponse {
  /// Returns a new [AuthResponse] instance.
  AuthResponse({
    required this.token,
    required this.expiresIn,
  });

  /// JWT access token
  String token;

  /// Token TTL in seconds (handler currently hard-codes 3600)
  int expiresIn;

  @override
  bool operator ==(Object other) => identical(this, other) || other is AuthResponse &&
    other.token == token &&
    other.expiresIn == expiresIn;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (token.hashCode) +
    (expiresIn.hashCode);

  @override
  String toString() => 'AuthResponse[token=$token, expiresIn=$expiresIn]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'token'] = this.token;
      json[r'expires_in'] = this.expiresIn;
    return json;
  }

  /// Returns a new [AuthResponse] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static AuthResponse? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "AuthResponse[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "AuthResponse[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return AuthResponse(
        token: mapValueOfType<String>(json, r'token')!,
        expiresIn: mapValueOfType<int>(json, r'expires_in')!,
      );
    }
    return null;
  }

  static List<AuthResponse> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <AuthResponse>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = AuthResponse.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, AuthResponse> mapFromJson(dynamic json) {
    final map = <String, AuthResponse>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = AuthResponse.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of AuthResponse-objects as value to a dart map
  static Map<String, List<AuthResponse>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<AuthResponse>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = AuthResponse.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'token',
    'expires_in',
  };
}

