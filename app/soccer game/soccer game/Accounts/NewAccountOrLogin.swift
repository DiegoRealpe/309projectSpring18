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
import Alamofire


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
        
        let screenSize:CGRect = UIScreen.main.bounds
        let screenHeight = screenSize.height //real screen height
        //let's suppose we want to have 10 points bottom margin
        let newCenterY = screenHeight - loginButton.frame.height - 10
        let newCenter = CGPoint(x: view.center.x,y:  newCenterY)
     
        
        loginButton.center = newCenter//view.center
        
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
    
    func sendCRUDServiceStatsRequest(FBToken : String) {
        
        let requestURL = "http://\(CommunicationProperties.crudServiceHost):\(CommunicationProperties.crudServicePort)/player/stats"
        
        let headers: HTTPHeaders = [
            "FacebookToken": FBToken,
            ]
        
        print("\n")
        print("\n")
        
        print("logging in with CRUD Service at URL",requestURL,"\nwith headers : \(headers)")
        
        Alamofire.request(requestURL , method : .get , headers : headers)
            .responseString(completionHandler: statsRequestResponse(_:))
    }

    func statsRequestResponse(_ response : DataResponse<String>)
    {
        
        print("🍆🍆🍆status code",response.response!.statusCode) //todo handle 404 without it being a fatal error
        
        
        if(response.response!.statusCode == 202)
        {
          //  buildStatsNodes();
        }
        
        
    }
    
}

