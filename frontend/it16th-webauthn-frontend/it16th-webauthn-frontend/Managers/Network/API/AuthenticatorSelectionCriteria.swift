//
//  AuthenticatorSelectionCriteria.swift
//  it16th-webauthn-frontend
//
//  Created by Leo Ho on 2024/9/20.
//

import Foundation

struct AuthenticatorSelectionCriteria: Codable {
    
    var authenticatorAttachment: String
    
    var residentKey: String
    
    var requireResidentKey: Bool
    
    var userVerification: String
    
    init(authenticatorAttachment: AuthenticatorAttachment,
         residentKey: ResidentKeyRequirement,
         requireResidentKey: Bool = false,
         userVerification: UserVerificationRequirement = .preferred) {
        self.authenticatorAttachment = authenticatorAttachment.rawValue
        self.residentKey = residentKey.rawValue
        self.requireResidentKey = requireResidentKey
        self.userVerification = userVerification.rawValue
    }
}

enum AuthenticatorAttachment: String, Codable {
    
    case platform = "platform"
    
    case crossPlatform = "cross-platform"
}

enum ResidentKeyRequirement: String, Codable {
    
    case required = "required"
    
    case preferred = "preferred"
    
    case discouraged = "discouraged"
}

enum UserVerificationRequirement: String, Codable {
    
    case required = "required"
    
    case preferred = "preferred"
    
    case discouraged = "discouraged"
}
