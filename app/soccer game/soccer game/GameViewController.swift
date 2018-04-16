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

class GameViewController:


UIViewController,FBSDKLoginButtonDelegate {
    
    @IBOutlet weak var chatView: ChatView!
    static var globalChatView: ChatView!
    
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
               //sendCRUDServiceLoginRequest(FBToken : AccessToken.current.unsafelyUnwrapped.authenticationToken)
            }
            
            
           //sendCRUDServiceLoginRequest(FBToken : AccessToken.current!.authenticationToken)
            loginButton.removeFromSuperview()
        }
    }
    
    func loginButtonDidLogOut(_ loginButton: FBSDKLoginButton!) {
        let loginButton = FBSDKLoginButton()
        loginButton.center = view.center
        loginButton.delegate = self // Remember to set the delegate of the loginButton
        view.addSubview(loginButton)
    }
    
    func setGameScene(/*_ loginButton: FBSDKLoginButton!*/)
    {
        
        if let view = self.view as! SKView? {
            // Load the SKScene from 'GameScene.sks'
            if let scene = MainMenu(fileNamed: "MainMenu") {
                // Set the scale mode to scale to fit the window
                scene.scaleMode = .aspectFill
                
                // Present the scene
                view.presentScene(scene)
            }
           // loginButton.removeFromSuperview()
            view.ignoresSiblingOrder = true
            
            view.showsFPS = true
            view.showsNodeCount = true
            
        }
        
        hideChatView()
        GameViewController.globalChatView = chatView
        
        
    }
    
    

    override func viewDidLoad() {
        super.viewDidLoad()
        
        setGameScene()
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
    
    

    
    struct loginCompletion {
        var finished = false
        var token: Int
    }
    
    func hideChatView(){
        self.chatView.isHidden = true
    }

    
    func unhideChatView(){
        self.chatView.isHidden = false
    }
    //currently not doing anything until crud service is sorted out
    /*func loginResponse(_ response : DataResponse<String>){
        print("login response was",String(describing : response.data))
        
    }*/

}
