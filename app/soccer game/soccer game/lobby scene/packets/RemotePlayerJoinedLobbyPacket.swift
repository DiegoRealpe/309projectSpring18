//
//  RemotePlayerJoinedLobby.swift
//  soccer game
//
//  Created by rtoepfer on 4/5/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation

class RemotePlayerJoinedLobbyPacket {
    
    var playerNumber : Int
    var username : String
    
    init(data : [UInt8]){
        self.playerNumber = Int(data[1])
        
        let usernameBytes = Array(data[2..<82])
        self.username = String(bytes: usernameBytes, encoding: .utf8)!
        self.username = self.username.replacingOccurrences(of: "\0", with: "")
        
    }
    
}
