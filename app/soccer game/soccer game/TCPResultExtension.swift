//
//  TCPResultExtension.swift
//  soccer game
//
//  Created by rtoepfer on 2/12/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation
import SwiftSocket

extension Result{
    
    func logError() {
        logError(message : "there was a socket error:")
    }
    
    func logError(message : String) {
        if self.isFailure, let err = self.error{
            print(message)
            print(err.localizedDescription)
        }
    }
    
}

