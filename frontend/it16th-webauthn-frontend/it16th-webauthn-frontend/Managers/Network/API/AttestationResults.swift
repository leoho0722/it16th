//
//  AttestationResults.swift
//  it16th-webauthn-frontend
//
//  Created by Leo Ho on 2024/9/20.
//

import Foundation

struct AttestationResultsRequest: Codable {
    
    var id: String
    
    var response: AuthenticatorAttestationResponse
    
    var getClientExtensionResults: ClientExtensionResults
    
    var type: String
    
    struct AuthenticatorAttestationResponse: Codable {
        
        var clientDataJSON: String
        
        var attestationObject: String
    }
}