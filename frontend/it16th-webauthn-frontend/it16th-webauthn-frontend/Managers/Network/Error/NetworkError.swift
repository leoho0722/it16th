//
//  NetworkError.swift
//  it16th-webauthn-frontend
//
//  Created by Leo Ho on 2024/9/21.
//

import Foundation

enum NetworkError: Error {
    
    case badRequest(Data)
    
    case internalServerError(Data)
    
    case jsonEncodeFailed(Error)
    
    case jsonDecodeFailed(Error)
    
    case unknown(Error)
}
