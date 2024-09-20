//
//  AssertionResults.swift
//  it16th-webauthn-frontend
//
//  Created by Leo Ho on 2024/9/20.
//

import Foundation

struct AssertionResultsRequest: Codable {
    
    var id: String
    
    var response: AuthenticatorAssertionResponse
    
    var getClientExtensionResults: ClientExtensionResults
    
    var type: String
    
    struct AuthenticatorAssertionResponse: Codable {
        
        var authenticatorData: String
        
        var signature: String
        
        var userHandle: String
        
        var clientDataJSON: String
    }
}