//
//  NetworkManager.swift
//  it16th-webauthn-frontend
//
//  Created by Leo Ho on 2024/9/20.
//

import Foundation
import SwiftHelpers

actor NetworkManager: NSObject {

    static let shared = NetworkManager()
    
    private let urlSessionConfiguration: URLSessionConfiguration
    private let urlSession: URLSession
    
    override init() {
        self.urlSessionConfiguration = .default
        self.urlSession = URLSession(configuration: self.urlSessionConfiguration)
    }
    
    func request<D>(with config: RequestConfiguration) async throws -> D where D: Decodable {
        let request = try buildURLRequest(config: config)
        let (data, response) = try await urlSession.data(for: request)
        guard let httpResponse = (response as? HTTPURLResponse) else {
            throw URLError(.badServerResponse)
        }
        switch HTTP.StatusCode(rawValue: httpResponse.statusCode) {
        case .badRequest:
            throw NetworkError.badRequest(data)
        case .internalServerError:
            throw NetworkError.internalServerError(data)
        case .ok:
            let decodedResponse: D = try decodeResponse(data: data)
            return decodedResponse
        default:
            let decodedResponse: D = try decodeResponse(data: data)
            return decodedResponse
        }
    }
    
    private func buildURLRequest(config: RequestConfiguration) throws -> URLRequest {
        guard let url = URL(string: "\(config.scheme.rawValue)\(config.host.rawValue)\(config.endpoint.rawValue)") else {
            throw URLError(.badURL)
        }
        var request = URLRequest(url: url)
        request.httpMethod = config.method.rawValue
        request.allHTTPHeaderFields = config.headers
        
        switch config.method {
        case .post:
            do {
                request.httpBody = try JSON.toJsonData(data: config.body)
            } catch {
                throw NetworkError.jsonEncodeFailed(error)
            }
        default:
            break
        }
        
        return request
    }
    
    private func decodeResponse<D>(data: Data) throws -> D where D: Decodable {
        do {
            let decoder = JSONDecoder()
            let decodedResponse = try decoder.decode(D.self, from: data)
            return decodedResponse
        } catch {
            throw NetworkError.jsonDecodeFailed(error)
        }
    }
}
