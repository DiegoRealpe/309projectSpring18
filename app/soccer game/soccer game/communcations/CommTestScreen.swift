//
//  CommTestScreen.swift
//  soccer game
//
//  Created by rtoepfer on 2/9/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import SpriteKit
import Alamofire

class CommTestScreen: SKScene {
    
    private var label:SKLabelNode?
    
    override func didMove(to view: SKView) {
        
        print("at comm test screen")
        
        label = self.childNode(withName: "Debug Label") as? SKLabelNode
        
        Alamofire.request("http://proj-309-MG-6.cs.iastate.edu", method: .get).response { (response) in
            self.recievedResponse(response);
        }
        
    }
    
    func recievedResponse(_ response : DefaultDataResponse){
        
        if let data:Data = response.data{
            let str:String = String(data: data, encoding: .utf8)!
            print("got response " + str)
            
            self.label!.text = str
        }
    }
    
    
}

