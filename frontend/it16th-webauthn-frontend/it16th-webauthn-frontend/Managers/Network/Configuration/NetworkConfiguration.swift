//
//  NetworkConfiguration.swift
//  it16th-webauthn-frontend
//
//  Created by Leo Ho on 2024/9/20.
//

import Foundation
import SwiftHelpers

struct NetworkConfiguration {
    
    var method: HTTP.Method
    
    var scheme: NetworkScheme
    
    var host: String
    
    var endpoint: String
    
    var headers: Dictionary<String, String>
    
    var body: Encodable
    
    enum NetworkScheme: String {
        
        case http = "http://"
        
        case https = "https://"
    }
}
