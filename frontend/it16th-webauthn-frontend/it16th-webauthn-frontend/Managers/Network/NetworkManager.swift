//
//  NetworkManager.swift
//  it16th-webauthn-frontend
//
//  Created by Leo Ho on 2024/9/20.
//

import Foundation
import SwiftHelpers

class NetworkManager: NSObject {

    static let shared = NetworkManager()
    
    private let urlSessionConfiguration: URLSessionConfiguration
    private let urlSession: URLSession
    
    override init() {
        self.urlSessionConfiguration = .default
        self.urlSession = URLSession(configuration: self.urlSessionConfiguration)
    }
    
    func request<D>(with config: NetworkConfiguration) async throws -> Response<D> where D: Decodable {
        let request = try buildURLRequest(config: config)
        let (data, response) = try await urlSession.data(for: request)
        guard let httpResponse = (response as? HTTPURLResponse) else {
            throw URLError(.badServerResponse)
        }
        let decodedResponse: D = try decodeResponse(data: data)
        return Response(statusCode: httpResponse.statusCode, body: decodedResponse)
    }
    
    private func buildURLRequest(config: NetworkConfiguration) throws -> URLRequest {
        guard let url = URL(string: "\(config.scheme.rawValue)\(config.host)\(config.endpoint)") else {
            throw URLError(.badURL)
        }
        var request = URLRequest(url: url)
        request.httpMethod = config.method.rawValue
        request.allHTTPHeaderFields = config.headers
        
        switch config.method {
        case .post:
            request.httpBody = try JSON.toJsonData(data: config.body)
        default:
            break
        }
        
        return request
    }
    
    private func decodeResponse<D>(data: Data) throws -> D where D: Decodable {
        let decoder = JSONDecoder()
        return try decoder.decode(D.self, from: data)
    }
}
