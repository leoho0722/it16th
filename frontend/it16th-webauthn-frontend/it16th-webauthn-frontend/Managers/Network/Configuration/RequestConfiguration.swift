//
//  RequestConfiguration.swift
//  it16th-webauthn-frontend
//
//  Created by Leo Ho on 2024/9/20.
//

import Foundation
import SwiftHelpers

struct RequestConfiguration {
    
    var method: HTTP.Method
    
    var scheme: NetworkScheme
    
    var host: Host
    
    var endpoint: Endpoint
    
    var headers: Dictionary<String, String>
    
    var body: Encodable
    
    init(method: HTTP.Method,
         scheme: NetworkScheme,
         host: Host,
         endpoint: Endpoint,
         headers: Dictionary<String, String> = defaultHeaders(),
         body: Encodable) {
        self.method = method
        self.scheme = scheme
        self.host = host
        self.endpoint = endpoint
        self.headers = headers
        self.body = body
    }
}

extension RequestConfiguration {
    
    enum NetworkScheme: String {
        
        case http = "http://"
        
        case https = "https://"
    }
    
    enum Host: String {
        
        case rpServer = "ab15-211-20-7-115.ngrok-free.app"
    }
    
    enum Endpoint: String {
        
        case beginRegistration = "/attestation/options"
        
        case finishRegistration = "/attestation/result"
        
        case beginAuthentication = "/assertion/options"
        
        case finishAuthentication = "/assertion/result"
    }
    
    static func defaultHeaders() -> Dictionary<String, String> {
        return [
            HTTP.HeaderFields.contentType.rawValue : HTTP.ContentType.json.rawValue
        ]
    }
}
