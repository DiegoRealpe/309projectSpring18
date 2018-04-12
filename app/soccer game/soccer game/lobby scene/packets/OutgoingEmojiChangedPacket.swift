//
//  OutgoingEmojiChangedPacket.swift
//  soccer game
//
//  Created by rtoepfer on 4/12/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation

class OutgoingEmojiChangedPacket : SendablePacket{

    var emoji : String
    
    init(emoji : String){
        self.emoji = emoji
    }
    
    func toByteArray() -> [UInt8] {
        var array = [UInt8]()
        array.reserveCapacity(25)
        
        array.append(208)
        
        let utf8Array : [UInt8] = Array(emoji.utf8)
        array.append(contentsOf: utf8Array)
        array.append(contentsOf: Array(repeating: 0, count: 24-utf8Array.count))
        
        print(array, array.count)
        
        return array
    }
}
