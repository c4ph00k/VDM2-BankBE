//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class SignUpRequest {
  /// Returns a new [SignUpRequest] instance.
  SignUpRequest({
    required this.email,
    required this.password,
    required this.username,
    required this.firstName,
    required this.lastName,
    required this.fiscalCode,
  });

  String email;

  String password;

  String username;

  String firstName;

  String lastName;

  String fiscalCode;

  @override
  bool operator ==(Object other) => identical(this, other) || other is SignUpRequest &&
    other.email == email &&
    other.password == password &&
    other.username == username &&
    other.firstName == firstName &&
    other.lastName == lastName &&
    other.fiscalCode == fiscalCode;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (email.hashCode) +
    (password.hashCode) +
    (username.hashCode) +
    (firstName.hashCode) +
    (lastName.hashCode) +
    (fiscalCode.hashCode);

  @override
  String toString() => 'SignUpRequest[email=$email, password=$password, username=$username, firstName=$firstName, lastName=$lastName, fiscalCode=$fiscalCode]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'email'] = this.email;
      json[r'password'] = this.password;
      json[r'username'] = this.username;
      json[r'first_name'] = this.firstName;
      json[r'last_name'] = this.lastName;
      json[r'fiscal_code'] = this.fiscalCode;
    return json;
  }

  /// Returns a new [SignUpRequest] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static SignUpRequest? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "SignUpRequest[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "SignUpRequest[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return SignUpRequest(
        email: mapValueOfType<String>(json, r'email')!,
        password: mapValueOfType<String>(json, r'password')!,
        username: mapValueOfType<String>(json, r'username')!,
        firstName: mapValueOfType<String>(json, r'first_name')!,
        lastName: mapValueOfType<String>(json, r'last_name')!,
        fiscalCode: mapValueOfType<String>(json, r'fiscal_code')!,
      );
    }
    return null;
  }

  static List<SignUpRequest> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <SignUpRequest>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = SignUpRequest.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, SignUpRequest> mapFromJson(dynamic json) {
    final map = <String, SignUpRequest>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = SignUpRequest.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of SignUpRequest-objects as value to a dart map
  static Map<String, List<SignUpRequest>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<SignUpRequest>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = SignUpRequest.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'email',
    'password',
    'username',
    'first_name',
    'last_name',
    'fiscal_code',
  };
}

