//
//  NewAccountOrLogin.swift
//  soccer game
//
//  Created by Mark Schwartz on 2/13/18.
//  Copyright © 2018 MG 6. All rights reserved.
//
import SpriteKit
import UIKit
import FacebookCore
import FacebookLogin
import FBSDKLoginKit
//import Alamofire



class NewAccountOrLogin: SKScene{
    
    
    var back : SKNode?
    var logout: SKNode?
    var viewController: UIViewController?
    var logOutLabel: SKLabelNode?
    var isLoggedIn = AccessToken.current != nil
    var loginButton = LoginButton(readPermissions: [ .publicProfile ])
    
    override func didMove(to view: SKView) {
        print("got to accounts")
        
        //get scene subnodes
        self.back = self.childNode(withName: "Back Button")
        
        
     
        loginButton.center = view.center
        
        view.addSubview(loginButton)
        
    }
    
    override func touchesBegan(_ touches: Set<UITouch>, with event: UIEvent?)
    {
        
        
        //if necesarry nodes and one touch exist
        if let t = touches.first ,let back = self.back
        {
            let point = t.location(in: self)
            
            //see if touch contains first
            if (back.contains(point))
            {
                self.moveToScene(.mainMenu)
                loginButton.removeFromSuperview()
            }
          
                
                
    
        }
    }
        
    
    
    func fadeInLabel(label : SKLabelNode?){
        if let nonOptLabel = label{
            nonOptLabel.alpha = 0.0
            nonOptLabel.run(SKAction.fadeIn(withDuration: 2.0))
        }
    }
    

    
}

