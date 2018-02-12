//
//  CommTestScreen.swift
//  soccer game
//
//  Created by rtoepfer on 2/9/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import SpriteKit
import Alamofire

//for a screen to test network communications in an enviroment with no reprecussions, will not be in final producrt
class CommTestScreen: SKScene {
    
    //label for output
    private var label:SKLabelNode?
    
    //sandbox here
    override func didMove(to view: SKView) {
        
        print("at comm test screen")
        
        label = self.childNode(withName: "Debug Label") as? SKLabelNode
        
        testHttp() //use refer string response to recievedResponse
    }
    
    fileprivate func testHttp(){
        Alamofire.request("http://localhost:8080", method: .get)
                .responseString(completionHandler: recievedResponse(_:))
    }
    
    //to be passed to alamofire to handle string response
    //needs future multithreading support
    func recievedResponse(_ response : DataResponse<String>){
        
        if let data:Data = response.data, let str:String = String(data: data, encoding: .utf8){
            print("got response " + str)
            
            self.label!.text = str
        }
    }
    
    
}

