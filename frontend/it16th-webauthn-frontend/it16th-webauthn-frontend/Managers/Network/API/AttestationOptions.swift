//
//  AttestationOptions.swift
//  it16th-webauthn-frontend
//
//  Created by Leo Ho on 2024/9/20.
//

import Foundation

struct AttestationOptionsRequest: Codable {
    
    var username: String
    
    var displayName: String
    
    var authenticatorSelection: AuthenticatorSelectionCriteria
    
    var attestation: String
}

struct AttestationOptionsResponse: Decodable {
    
    let status: String
    
    let errorMessage: String
    
    let rp: RelyingParty
    
    let user: UserEntity
    
    let challenge: String
    
    let pubKeyCredParams: [PubKeyCredParam]
    
    let timeout: Int
    
    let excludeCredentials: [ExcludeCredential]
    
    let authenticatorSelection: AuthenticatorSelectionCriteria
    
    let attestation: String
    
    struct RelyingParty: Decodable {
        
        let name: String
    }
    
    struct UserEntity: Decodable {
        
        let id: String
        
        let name: String
        
        let displayName: String
    }
    
    struct PubKeyCredParam: Decodable {
        
        let type: String
        
        let alg: Int
    }
    
    struct ExcludeCredential: Decodable {
        
        let type: String
        
        let id: String
    }
}
