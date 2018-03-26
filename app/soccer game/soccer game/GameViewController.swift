//
//  GameViewController.swift
//  soccer game
//
//  Created by Mark Schwartz on 1/28/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import UIKit
import SpriteKit
import GameplayKit
import FBSDKLoginKit
import FacebookCore
import FacebookLogin
import Alamofire

class GameViewController: UIViewController,FBSDKLoginButtonDelegate {
    
    
    func loginButton(_ loginButton: FBSDKLoginButton!, didCompleteWith result: FBSDKLoginManagerLoginResult!, error: Error!)
    {
        
        if ((error) != nil)
        {
            // Process error
        }
        else if result.isCancelled {
            // Handle cancellations
        }
        else {
            // If you ask for multiple permissions at once, you
            // should check if specific permissions missing
            if result.grantedPermissions.contains("email")
            {
               
            }
            print("LoggedInThroughFB")
            
            sendCRUDServiceLoginRequest(FBToken : AccessToken.current!.authenticationToken)
        }
    }
    
    func loginButtonDidLogOut(_ loginButton: FBSDKLoginButton!) {
        let loginButton = FBSDKLoginButton()
        loginButton.center = view.center
        loginButton.delegate = self // Remember to set the delegate of the loginButton
        view.addSubview(loginButton)
    }
    
    func setGameScene(_ loginButton: FBSDKLoginButton!)
    {
        if let view = self.view as! SKView? {
            // Load the SKScene from 'GameScene.sks'
            if let scene = MainMenu(fileNamed: "MainMenu") {
                // Set the scale mode to scale to fit the window
                scene.scaleMode = .aspectFill
                
                // Present the scene
                view.presentScene(scene)
            }
            loginButton.removeFromSuperview()
            view.ignoresSiblingOrder = true
            
            view.showsFPS = true
            view.showsNodeCount = true
        }
        
        print("token was",AccessToken.current)
        sendCRUDServiceLoginRequest(FBToken: AccessToken.current!.authenticationToken)
    }
    
    

    override func viewDidLoad() {
        super.viewDidLoad()
        let loginButton = FBSDKLoginButton()
        // Do any additional setup after loading the view, typically from a nib.
        
        loginIfNecesarry()
        
    }

    override var shouldAutorotate: Bool {
        return true
    }

    override var supportedInterfaceOrientations: UIInterfaceOrientationMask {
        if UIDevice.current.userInterfaceIdiom == .phone {
            return .allButUpsideDown
        } else {
            return .all
        }
    }

    override func didReceiveMemoryWarning() {
        super.didReceiveMemoryWarning()
        // Release any cached data, images, etc that aren't in use.
    }

    override var prefersStatusBarHidden: Bool {
        return true
    }
    
    func loginIfNecesarry(){
        //todo: remove this once we are asking game service for login
        if let accessToken = AccessToken.current
        {
            sendCRUDServiceLoginRequest(FBToken: accessToken.authenticationToken)
        }
            
        else
        {
            let loginButton = FBSDKLoginButton()
            loginButton.center = view.center
            loginButton.delegate = self // Remember to set the delegate of the loginButton
            view.addSubview(loginButton)
        }
        
    }
    
    
    //returns application token
    func sendCRUDServiceLoginRequest(FBToken : String) {
        
        let requestURL = "\(CommunicationProperties.crudServiceHost):\(CommunicationProperties.crudServicePort)/player/login"
        
        let headers: HTTPHeaders = [
            "FacebookToken": FBToken,
            ]
        
        print("logging in with CRUD Service at URL",requestURL,"\nwith headers : \(headers)")
        
        Alamofire.request(requestURL , method : .get , headers : headers)
            .responseString(completionHandler: loginResponse(_:))
    }
    
    struct loginCompletion {
        var finished = false
        var token: Int
    }
    
    //currently not doing anything until crud service is sorted out
    func loginResponse(_ response : DataResponse<String>){
        print("login response was",String(describing : response.data))
        
    }

}
