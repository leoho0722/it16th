//
//  Common.swift
//  it16th-webauthn-frontend
//
//  Created by Leo Ho on 2024/9/20.
//

import Foundation

struct CommonResponse: Decodable {
    
    let status: String
    
    let errorMessage: String
}
