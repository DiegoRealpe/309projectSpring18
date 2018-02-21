//
//  SKSceneTransitionExtension.swift
//  soccer game
//
//  Created by rtoepfer on 1/29/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import SpriteKit

enum Scene {
    case gameScene
    case mainMenu
    case matchMakingScene
    case commTestScreen
}

let sceneToFileDict : [Scene:String] = [
    .gameScene : "GameScene",
    .mainMenu : "MainMenu",
    .matchMakingScene : "MatchMakingScene",
    .commTestScreen : "CommTestScreen"
]

extension SKScene{
    
    func moveToScene(_ scene : Scene){
        moveToScene(scene, dataFunction : { (dict) in })
    }
    
    func moveToScene(_ scene : Scene, dataFunction : (NSMutableDictionary) -> Void){
        let sceneName = sceneToFileDict[scene]!
        
        if let skScene = GameScene(fileNamed: sceneName){
            // Set the scale mode to scale to fit the window
            skScene.scaleMode = .aspectFill
            skScene.userData = NSMutableDictionary()
            let userData: NSMutableDictionary = skScene.userData!
            dataFunction(userData)
            
            // Present the scene
            if let view = self.view{
                view.presentScene(skScene)
                
            }
        }
    }
    
   
    func moveToAccountScreen(){
        
        if let scene = GameScene(fileNamed: "NewAccountOrLogin"){
            // Set the scale mode to scale to fit the window
            scene.scaleMode = .aspectFill
            
            
            // Present the scene
            if let view = self.view{
                view.presentScene(scene)
            }
        }
        
    }
    
   
    
    
}
