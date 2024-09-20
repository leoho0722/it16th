//
//  AssertionOptions.swift
//  it16th-webauthn-frontend
//
//  Created by Leo Ho on 2024/9/20.
//

import Foundation

struct AssertionOptionsRequest: Codable {
    
    var username: String
    
    var userVerification: String
    
    init(username: String, userVerification: UserVerificationRequirement) {
        self.username = username
        self.userVerification = userVerification.rawValue
    }
}

struct AssertionOptionsResponse: Decodable {
    
    let status: String
    
    let errorMessage: String
    
    let challenge: String
    
    let timeout: Int
    
    let rpId: String
    
    let allowCredentials: [AllowCredential]
    
    let userVerification: String
    
    struct AllowCredential: Decodable {
        
        let id: String
        
        let type: String
    }
}
