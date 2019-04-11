#include <string>
#include <array>
#include <fstream>
#include <iostream>
#include <iterator>
#include <algorithm>
#include <vector>
#include <utility>
#include <tuple>
#include <functional>

using namespace std;
using namespace std::placeholders;

pair<unsigned int, string> countGetRef(const string f);
void outputSNPInfo(string file, string& refSeq, const unsigned refLen);

int main(int argc, char* argv[])
{ 
  if(argc == 1)
  {
    //printHelp();
    //cout << ((('A'-'a')%32==0) ? "true" : "false") << endl;
    //cout << ('a'-65)%32 << endl;
    //cout << ('A'-'T')%32 << endl;
    return 0;
  }

  // get file list
  vector<string> fileList;
  for(size_t i = 1; i < argc; ++i)
  {
    fileList.push_back(argv[i]);
  }

  // Determine mask length then get ref seq.
  string refSeq;
  unsigned int refLen, minLen;
  std::tie(refLen, refSeq) = countGetRef(fileList[0]);

  //for_each(fileList.begin(), fileList.end(), setMaskLambda);
  for(auto it=fileList.begin(); it !=fileList.end(); ++it)
  {
    outputSNPInfo(*it, refSeq, refLen);
  }

  return 0;
}

pair<unsigned int, string> countGetRef(const string f)
{
  ifstream ifs(f);
  unsigned int count(0);
  string s;
  char c;

  ifs.seekg(0, ifs.end);
  string ref(ifs.tellg(), 'n');
  ifs.seekg(0, ifs.beg);

  getline(ifs, s);
  while(ifs.get(c))
  {
    if(c != 10){
      ref[count] = c;
      count++;
    }
  }
  ifs.close();

  ref.resize(count);
  return make_pair( count, ref);
}
void outputSNPInfo(string file, string& refSeq, const unsigned refLen)
{
  string::size_type start, length;
  start = file.rfind("/") + 1;
  length = file.rfind(".") - start;
  cout << file.substr(start, length) << ",";
  ifstream ifs(file);
  unsigned int count(0);
  string s;
  char c;

  getline(ifs, s);
  while(ifs.get(c))
  {
    if(c != 10)
    {
      if((refSeq[count]-c) % 32 == 0){
        cout <<  "1,";
      }else{
        if (c == 'n' || c == 'N'){
          cout <<  "0,";
        }else{
          cout <<  "2,";
        }
      }
      ++count;

      if(count >= refLen)
        break;
    }
  }
  cout << '\n' ;
  ifs.close();
}