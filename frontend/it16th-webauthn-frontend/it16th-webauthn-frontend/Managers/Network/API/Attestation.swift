//
//  Attestation.swift
//  it16th-webauthn-frontend
//
//  Created by Leo Ho on 2024/9/21.
//

import Foundation

enum Attestation: String, Codable {
    
    case none = "none"
    
    case direct = "direct"
    
    case indirect = "indirect"
}
