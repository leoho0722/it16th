//
//  PaskeysManager.swift
//  it16th-webauthn-frontend
//
//  Created by Leo Ho on 2024/9/19.
//

import AuthenticationServices
import Foundation
import os

class PasskeysManager: NSObject {
    
    private let domain: String
    
    var authenticationAnchor: ASPresentationAnchor?
    
    weak var delegate: PasskeysManagerDelegate?
    
    private let logger = Logger()
    
    init(domain: String) {
        self.domain = domain
    }
    
    // MARK: PassKeys Registration
    
    func registration(username: String, challenge: String, anchor: ASPresentationAnchor) {
        
        self.authenticationAnchor = anchor
        
        let publicKeyCredentialProvider = ASAuthorizationPlatformPublicKeyCredentialProvider(relyingPartyIdentifier: domain)
        
        // Fetch the challenge from the server. The challenge needs to be unique for each request.
        // The userID is the identifier for the user's account.
        let challenge = Data(challenge.utf8)
        let userID = Data(username.utf8)
        
        let registrationRequest = publicKeyCredentialProvider.createCredentialRegistrationRequest(challenge: challenge,
                                                                                                  name: username,
                                                                                                  userID: userID)
        
        // Use only ASAuthorizationPlatformPublicKeyCredentialRegistrationRequests or
        // ASAuthorizationSecurityKeyPublicKeyCredentialRegistrationRequests here.
        let authController = ASAuthorizationController(authorizationRequests: [registrationRequest])
        authController.delegate = self
        authController.presentationContextProvider = self
        authController.performRequests()
    }
    
    // MARK: PassKeys Authentication
    
    func authentication(challenge: String,
                        anchor: ASPresentationAnchor,
                        preferImmediatelyAvailableCredentials: Bool) {
        
        self.authenticationAnchor = anchor
        
        let publicKeyCredentialProvider = ASAuthorizationPlatformPublicKeyCredentialProvider(relyingPartyIdentifier: domain)
        
        // Fetch the challenge from the server. The challenge needs to be unique for each request.
        let challenge = challenge.data(using: .utf8)!
        
        let assertionRequest = publicKeyCredentialProvider.createCredentialAssertionRequest(challenge: challenge)
        
        // Pass in any mix of supported sign-in request types.
        let authController = ASAuthorizationController(authorizationRequests: [assertionRequest])
        authController.delegate = self
        authController.presentationContextProvider = self
        
        if preferImmediatelyAvailableCredentials {
            // If credentials are available, presents a modal sign-in sheet.
            // If there are no locally saved credentials, no UI appears and
            // the system passes ASAuthorizationError.Code.canceled to call
            // `PasskeysManager.authorizationController(controller:didCompleteWithError:)`.
            authController.performRequests(options: .preferImmediatelyAvailableCredentials)
        } else {
            // If credentials are available, presents a modal sign-in sheet.
            // If there are no locally saved credentials, the system presents a QR code to allow signing in with a
            // passkey from a nearby device.
            authController.performRequests()
        }
    }
    
    // MARK: PassKeys Authentication with AutoFill
    
    func authenticationWithAutoFill(challenge: String, anchor: ASPresentationAnchor) {
        
        self.authenticationAnchor = anchor
        
        let publicKeyCredentialProvider = ASAuthorizationPlatformPublicKeyCredentialProvider(relyingPartyIdentifier: domain)
        
        // Fetch the challenge from the server. The challenge needs to be unique for each request.
        let challenge = Data(challenge.utf8)
        let assertionRequest = publicKeyCredentialProvider.createCredentialAssertionRequest(challenge: challenge)
        
        // AutoFill-assisted requests only support ASAuthorizationPlatformPublicKeyCredentialAssertionRequest.
        let authController = ASAuthorizationController(authorizationRequests: [assertionRequest])
        authController.delegate = self
        authController.presentationContextProvider = self
        authController.performAutoFillAssistedRequests()
    }
    
    // MARK: Passkeys Reset or Change
    
    func resetOrChangeWith(username: String, challenge: String, anchor: ASPresentationAnchor) {
        self.authenticationAnchor = anchor
        
        let publicKeyCredentialProvider = ASAuthorizationSecurityKeyPublicKeyCredentialProvider(relyingPartyIdentifier: domain)
        
        // Fetch the challenge from the server. The challenge needs to be unique for each request.
        let challenge = Data(challenge.utf8)
        let userID = Data(username.utf8)
        let registrationRequest = publicKeyCredentialProvider.createCredentialRegistrationRequest(challenge: challenge,
                                                                                                  displayName: username,
                                                                                                  name: username,
                                                                                                  userID: userID)
        // Use only ASAuthorizationPlatformPublicKeyCredentialRegistrationRequests or
        // ASAuthorizationSecurityKeyPublicKeyCredentialRegistrationRequests here.
        let authController = ASAuthorizationController(authorizationRequests: [registrationRequest])
        authController.delegate = self
        authController.presentationContextProvider = self
        authController.performRequests()
    }
}

extension PasskeysManager: ASAuthorizationControllerDelegate {
    
    func authorizationController(controller: ASAuthorizationController,
                                 didCompleteWithAuthorization authorization: ASAuthorization) {
        switch authorization.credential {
        case let credentialRegistration as ASAuthorizationPlatformPublicKeyCredentialRegistration:
            logger.log("A new passkey was registered: \(credentialRegistration)")
            // Verify the attestationObject and clientDataJSON with your service.
            // The attestationObject contains the user's new public key to store and use for subsequent sign-ins.
            
            (delegate as! PasskeysRegistration).passkeysManager(with: credentialRegistration)
            // After the server verifies the registration and creates the user account, sign in the user with the new account.
        case let credentialAssertion as ASAuthorizationPlatformPublicKeyCredentialAssertion:
            logger.log("A passkey was used to sign in: \(credentialAssertion)")
            // Verify the below signature and clientDataJSON with your service for the given userID.
            
            (delegate as! PasskeysAuthentication).passkeysManager(with: credentialAssertion)
            // After the server verifies the assertion, sign in the user.
        default:
            fatalError("Received unknown authorization type.")
        }
    }
    
    func authorizationController(controller: ASAuthorizationController,
                                 didCompleteWithError error: Error) {
        guard let authorizationError = error as? ASAuthorizationError else {
            logger.error("Unexpected authorization error: \(error.localizedDescription)")
            return
        }
        
        delegate?.passkeysManager(controller: controller,
                                  didCompleteWithError: authorizationError)
    }
}

extension PasskeysManager: ASAuthorizationControllerPresentationContextProviding {
    
    func presentationAnchor(for controller: ASAuthorizationController) -> ASPresentationAnchor {
        return authenticationAnchor!
    }
}

protocol PasskeysManagerDelegate: NSObjectProtocol {
    
    func passkeysManager(controller: ASAuthorizationController, didCompleteWithError error: Error)
}

protocol PasskeysRegistration: PasskeysManagerDelegate {
    
    func passkeysManager(with credentialRegistration: ASAuthorizationPlatformPublicKeyCredentialRegistration)
}

protocol PasskeysAuthentication: PasskeysManagerDelegate {
    
    func passkeysManager(with credentialAssertion: ASAuthorizationPlatformPublicKeyCredentialAssertion)
}
